"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const validator = require("validator");
const Account = require("../models/Account");
const ResetToken = require("../models/ResetToken");
const uuid4 = require("uuid4");
const random = require("randomstring");
const bcrypt = require("bcrypt");
const sha1 = require("sha1");
const moment = require("moment");
/**
 * Sets flash error message for a field
 *
 * @param {*} req Request object
 * @param {string} msg The error message
 * @param {string} field The field name
 * @param {*} body Response body
 */
function setFlashError(req, msg, field, body) {
    req.flash("error", msg);
    req.flash("field", field);
    req.flash("form", body);
}
/**
 * Validates res.body containing account details
 * @export
 * @param {Account} account
 * @returns {ValidationError}
 */
function validateAccount(account) {
    if (validator.isEmpty(account.first_name))
        return { msg: "First name is required", field: "first_name" };
    if (validator.isEmpty(account.last_name))
        return { msg: "Last name is required", field: "last_name" };
    if (validator.isEmpty(account.email))
        return { msg: "Email is required", field: "email" };
    if (!validator.isEmail(account.email))
        return { msg: "Email does not appear to be valid", field: "email" };
    if (validator.isEmpty(account.password))
        return { msg: "Password is required", field: "password" };
    if (!validator.isLength(account.password, { min: 8, max: 120 }))
        return { msg: "Please enter a password between 8 and 120 characters, the trickier the better!", field: "password" };
    return null;
}
exports.validateAccount = validateAccount;
/**
 * Handle signup processing
 *
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
function signup(req, res) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            let body = req.body;
            const valError = validateAccount(body);
            // send validation error
            if (valError != null) {
                setFlashError(req, valError.msg, valError.field, body);
                return res.redirect('back');
            }
            // check if an account with matching email exists
            let existingAccount = yield Account.get({ email: body.email });
            if (existingAccount) {
                setFlashError(req, "Email has already been registered", "email", body);
                return res.redirect('back');
            }
            // set id, client id and secret
            body.id = uuid4();
            body.password = bcrypt.hashSync(body.password, 10);
            body.confirmation_code = random.generate(28);
            let account = yield Account.create(body);
            // send confirmation message
            let resp = yield EmailService.sendAccountConfirmation(body.email, account.toJSON());
            req.flash("account", account);
            res.redirect("/confirm");
        }
        catch (e) {
            console.error(e);
            return res.serverError();
        }
    });
}
exports.signup = signup;
/**
 * Login processing
 *
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
function login(req, res) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            let body = req.body;
            if (!body.email) {
                setFlashError(req, 'Please enter your email', 'email', body);
                return res.redirect('back');
            }
            if (!body.password) {
                setFlashError(req, 'Please enter your password', 'password', body);
                return res.redirect('back');
            }
            let account = yield Account.get({ email: body.email });
            if (!account) {
                setFlashError(req, 'Your email and password combination is not valid', '', body);
                return res.redirect('back');
            }
            if (!bcrypt.compareSync(body.password, account.password)) {
                setFlashError(req, 'Your email and password combination is not valid', '', body);
                return res.redirect('back');
            }
            if (!account.confirmed) {
                req.flash('warn', 'warn');
                setFlashError(req, `Your account has not been confirmed. Please check your email for the confirmation message. <a href="/account/resend-code/${account.id}">Resend</a>`, '', body);
                return res.redirect('back');
            }
            return res.send(account);
        }
        catch (e) {
            return res.serverError();
        }
    });
}
exports.login = login;
/**
 * Verify an account
 *
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
function resendConfirmation(req, res) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            let userID = req.param("id");
            if (!userID) {
                return res.serverError();
            }
            // get account by user id
            let account = yield Account.get({ id: userID });
            if (!account) {
                req.flash('error', 'Sorry, This confirmation link is invalid');
                return res.view('account/verify');
            }
            if (account.confirmed) {
                req.flash('msg', 'This account has already been confirmed');
                return res.view('account/verify');
            }
            // send confirmation message
            let resp = yield EmailService.sendAccountConfirmation(account.email, account.toJSON());
            req.flash('msg', 'Confirmation message has been sent. Please check your email');
            return res.view('account/verify');
        }
        catch (e) {
            return res.serverError();
        }
    });
}
exports.resendConfirmation = resendConfirmation;
/**
 * Send reset password token to email
 *
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
function forgotPassword(req, res) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            let email = req.body.email;
            if (!email) {
                setFlashError(req, "Please enter your email address", "email", req.body);
                return res.view('account/forgot_password');
            }
            // get account by email
            let account = yield Account.get({ email });
            if (!account) {
                req.flash('error', `No account matching the email '${req.body.email}'`);
                return res.view('account/forgot_password');
            }
            let resetToken = yield ResetToken.create({ token: sha1(uuid4()), account: account.id });
            // send password reset message
            account = account.toJSON();
            account.token = resetToken.token;
            let resp = yield EmailService.sendPasswordResetEmail(account.email, account);
            req.flash('msg', 'A reset message has been sent to you');
            return res.view('account/forgot_password');
        }
        catch (e) {
            return res.serverError();
        }
    });
}
exports.forgotPassword = forgotPassword;
/**
 * Reset a password page
 *
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
function resetPasswordPage(req, res) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            let resetToken = yield ResetToken.get({ token: req.param("token") });
            if (!resetToken) {
                req.flash('error2', `Reset link is invalid or has expired`);
                return res.view('account/reset_password');
            }
            if (resetToken.used) {
                req.flash('error2', `Reset link is invalid or has expired`);
                return res.view('account/reset_password');
            }
            if (!moment(resetToken.created_at).add(2, "day").utc().isAfter(moment())) {
                req.flash('error2', `Reset link is invalid or has expired`);
            }
            return res.view('account/reset_password');
        }
        catch (e) {
            return res.serverError();
        }
    });
}
exports.resetPasswordPage = resetPasswordPage;
/**
 * Reset a password
 *
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
function resetPassword(req, res) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            let resetToken = yield ResetToken.get({ token: req.param("token") });
            if (!resetToken) {
                req.flash('error2', `Reset link is invalid or has expired`);
                return res.redirect(req.originalUrl);
            }
            if (resetToken.used) {
                req.flash('error2', `Reset link is invalid or has expired`);
                return res.redirect(req.originalUrl);
            }
            if (!moment(resetToken.created_at).add(2, "day").utc().isAfter(moment())) {
                req.flash('error2', `Reset link is invalid or has expired`);
                return res.redirect(req.originalUrl);
            }
            let body = req.body;
            // validation
            if (validator.isEmpty(body.password)) {
                setFlashError(req, `Password is required`, 'password', body);
                return res.redirect(req.originalUrl);
            }
            else if (!validator.isLength(body.password, { min: 8, max: 120 })) {
                setFlashError(req, `Please enter a password between 8 and 120 characters, the trickier the better!`, 'password', body);
                return res.redirect(req.originalUrl);
            }
            let account = yield Account.get({ id: resetToken.account });
            if (!account) {
                req.flash('error2', `Reset link is invalid or has expired'`);
                return res.redirect(req.originalUrl);
            }
            // ensure new password is not the same as the old one
            if (bcrypt.compareSync(body.password, account.password)) {
                setFlashError(req, `New password cannot be similar to a previously used password`, 'password', body);
                return res.redirect(req.originalUrl);
            }
            // set reset token as used
            resetToken.set('used', true);
            yield resetToken.save();
            // update password
            let newPasswordHash = bcrypt.hashSync(body.password, 10);
            account.set('password', newPasswordHash);
            yield account.save();
            return res.view('account/reset_password', { msg: 'Password successfully changed' });
        }
        catch (e) {
            return res.serverError();
        }
    });
}
exports.resetPassword = resetPassword;
/**
 * Verify an account
 *
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
function verify(req, res) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            let code = req.param("code");
            if (!code) {
                return res.serverError();
            }
            // get account by confirmation code
            let account = yield Account.get({ confirmation_code: code });
            if (!account) {
                req.flash('error', 'Sorry, This confirmation link is invalid');
                return res.view('account/verify');
            }
            if (account.confirmed) {
                req.flash('msg', 'This account has already been confirmed');
                return res.view('account/verify');
            }
            account.set('confirmed', true);
            yield account.save();
            req.flash('msg', 'Great! Your account has been confirmed');
            return res.view('account/verify');
        }
        catch (e) {
            console.error(e);
            return res.serverError();
        }
    });
}
exports.verify = verify;

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9jb250cm9sbGVycy9BY2NvdW50Q29udHJvbGxlci50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7O0FBR0EsdUNBQXdDO0FBQ3hDLDZDQUE4QztBQUM5QyxtREFBb0Q7QUFDcEQsK0JBQWdDO0FBQ2hDLHVDQUF3QztBQUV4QyxpQ0FBa0M7QUFDbEMsNkJBQThCO0FBQzlCLGlDQUFrQztBQUVsQzs7Ozs7OztHQU9HO0FBQ0gsdUJBQXVCLEdBQVEsRUFBRSxHQUFXLEVBQUUsS0FBYSxFQUFFLElBQVM7SUFDbEUsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLEVBQUUsR0FBRyxDQUFDLENBQUE7SUFDdkIsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLEVBQUUsS0FBSyxDQUFDLENBQUE7SUFDekIsR0FBRyxDQUFDLEtBQUssQ0FBQyxNQUFNLEVBQUUsSUFBSSxDQUFDLENBQUE7QUFDM0IsQ0FBQztBQUdEOzs7OztHQUtHO0FBQ0gseUJBQWdDLE9BQXVCO0lBQ25ELEVBQUUsQ0FBQyxDQUFDLFNBQVMsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLFVBQVUsQ0FBQyxDQUFDO1FBQUMsTUFBTSxDQUFDLEVBQUUsR0FBRyxFQUFFLHdCQUF3QixFQUFFLEtBQUssRUFBRSxZQUFZLEVBQUUsQ0FBQztJQUN6RyxFQUFFLENBQUMsQ0FBQyxTQUFTLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxTQUFTLENBQUMsQ0FBQztRQUFDLE1BQU0sQ0FBQyxFQUFFLEdBQUcsRUFBRSx1QkFBdUIsRUFBRSxLQUFLLEVBQUUsV0FBVyxFQUFFLENBQUM7SUFDdEcsRUFBRSxDQUFDLENBQUMsU0FBUyxDQUFDLE9BQU8sQ0FBQyxPQUFPLENBQUMsS0FBSyxDQUFDLENBQUM7UUFBQyxNQUFNLENBQUMsRUFBRSxHQUFHLEVBQUUsbUJBQW1CLEVBQUUsS0FBSyxFQUFFLE9BQU8sRUFBRSxDQUFDO0lBQzFGLEVBQUUsQ0FBQyxDQUFDLENBQUMsU0FBUyxDQUFDLE9BQU8s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QyxDQUFBO2dCQUMzRCxNQUFNLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxnQkFBZ0IsQ0FBQyxDQUFBO1lBQ3JDLENBQUM7WUFFRCxPQUFPLENBQUMsR0FBRyxDQUFDLFdBQVcsRUFBRSxJQUFJLENBQUMsQ0FBQTtZQUM5QixNQUFNLE9BQU8sQ0FBQyxJQUFJLEVBQUUsQ0FBQTtZQUNwQixHQUFHLENBQUMsS0FBSyxDQUFDLEtBQUssRUFBRSx3Q0FBd0MsQ0FBQyxDQUFBO1lBQzFELE1BQU0sQ0FBQyxHQUFHLENBQUMsSUFBSSxDQUFDLGdCQUFnQixDQUFDLENBQUE7UUFFckMsQ0FBQztRQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUM7WUFDVCxPQUFPLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ2hCLE1BQU0sQ0FBQyxHQUFHLENBQUMsV0FBVyxFQUFFLENBQUE7UUFDNUIsQ0FBQztJQUNMLENBQUM7Q0FBQTtBQTdCRCx3QkE2QkMiLCJmaWxlIjoiYXBpL2NvbnRyb2xsZXJzL0FjY291bnRDb250cm9sbGVyLmpzIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgc2V0RmxhZ3NGcm9tU3RyaW5nIH0gZnJvbSAndjgnO1xuaW1wb3J0IHsgYWNjZXNzIH0gZnJvbSAnZnMnO1xuaW1wb3J0IGV4cHJlc3MgPSByZXF1aXJlKCdleHByZXNzJyk7XG5pbXBvcnQgdmFsaWRhdG9yID0gcmVxdWlyZSgndmFsaWRhdG9yJyk7XG5pbXBvcnQgQWNjb3VudCA9IHJlcXVpcmUoJy4uL21vZGVscy9BY2NvdW50Jyk7XG5pbXBvcnQgUmVzZXRUb2tlbiA9IHJlcXVpcmUoJy4uL21vZGVscy9SZXNldFRva2VuJyk7XG5pbXBvcnQgdXVpZDQgPSByZXF1aXJlKCd1dWlkNCcpO1xuaW1wb3J0IHJhbmRvbSA9IHJlcXVpcmUoJ3JhbmRvbXN0cmluZycpO1xuaW1wb3J0IFJhbmRvbSA9IHJlcXVpcmUoJ3JhbmRvbS1qcycpO1xuaW1wb3J0IGJjcnlwdCA9IHJlcXVpcmUoJ2JjcnlwdCcpO1xuaW1wb3J0IHNoYTEgPSByZXF1aXJlKCdzaGExJyk7XG5pbXBvcnQgbW9tZW50ID0gcmVxdWlyZSgnbW9tZW50Jyk7XG5cbi8qKlxuICogU2V0cyBmbGFzaCBlcnJvciBtZXNzYWdlIGZvciBhIGZpZWxkXG4gKiBcbiAqIEBwYXJhbSB7Kn0gcmVxIFJlcXVlc3Qgb2JqZWN0XG4gKiBAcGFyYW0ge3N0cmluZ30gbXNnIFRoZSBlcnJvciBtZXNzYWdlXG4gKiBAcGFyYW0ge3N0cmluZ30gZmllbGQgVGhlIGZpZWxkIG5hbWVcbiAqIEBwYXJhbSB7Kn0gYm9keSBSZXNwb25zZSBib2R5XG4gKi9cbmZ1bmN0aW9uIHNldEZsYXNoRXJyb3IocmVxOiBhbnksIG1zZzogc3RyaW5nLCBmaWVsZDogc3RyaW5nLCBib2R5OiBhbnkpIHtcbiAgICByZXEuZmxhc2goXCJlcnJvclwiLCBtc2cpXG4gICAgcmVxLmZsYXNoKFwiZmllbGRcIiwgZmllbGQpXG4gICAgcmVxLmZsYXNoKFwiZm9ybVwiLCBib2R5KVxufVxuXG5cbi8qKlxuICogVmFsaWRhdGVzIHJlcy5ib2R5IGNvbnRhaW5pbmcgYWNjb3VudCBkZXRhaWxzXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge0FjY291bnR9IGFjY291bnQgXG4gKiBAcmV0dXJucyB7VmFsaWRhdGlvbkVycm9yfSBcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIHZhbGlkYXRlQWNjb3VudChhY2NvdW50OiBtb2RlbHMuQWNjb3VudCk6IFZhbGlkYXRpb25FcnJvciB7XG4gICAgaWYgKHZhbGlkYXRvci5pc0VtcHR5KGFjY291bnQuZmlyc3RfbmFtZSkpIHJldHVybiB7IG1zZzogXCJGaXJzdCBuYW1lIGlzIHJlcXVpcmVkXCIsIGZpZWxkOiBcImZpcnN0X25hbWVcIiB9O1xuICAgIGlmICh2YWxpZGF0b3IuaXNFbXB0eShhY2NvdW50Lmxhc3RfbmFtZSkpIHJldHVybiB7IG1zZzogXCJMYXN0IG5hbWUgaXMgcmVxdWlyZWRcIiwgZmllbGQ6IFwibGFzdF9uYW1lXCIgfTtcbiAgICBpZiAodmFsaWRhdG9yLmlzRW1wdHkoYWNjb3VudC5lbWFpbCkpIHJldHVybiB7IG1zZzogXCJFbWFpbCBpcyByZXF1aXJlZFwiLCBmaWVsZDogXCJlbWFpbFwiIH07XG4gICAgaWYgKCF2YWxpZGF0b3IuaXNFbWFpbChhY2NvdW50LmVtYWlsKSkgcmV0dXJuIHsgbXNnOiBcIkVtYWlsIGRvZXMgbm90IGFwcGVhciB0byBiZSB2YWxpZFwiLCBmaWVsZDogXCJlbWFpbFwiIH07XG4gICAgaWYgKHZhbGlkYXRvci5pc0VtcHR5KGFjY291bnQucGFzc3dvcmQpKSByZXR1cm4geyBtc2c6IFwiUGFzc3dvcmQgaXMgcmVxdWlyZWRcIiwgZmllbGQ6IFwicGFzc3dvcmRcIiB9O1xuICAgIGlmICghdmFsaWRhdG9yLmlzTGVuZ3RoKGFjY291bnQucGFzc3dvcmQsIHsgbWluOiA4LCBtYXg6IDEyMCB9KSkgcmV0dXJuIHsgbXNnOiBcIlBsZWFzZSBlbnRlciBhIHBhc3N3b3JkIGJldHdlZW4gOCBhbmQgMTIwIGNoYXJhY3RlcnMsIHRoZSB0cmlja2llciB0aGUgYmV0dGVyIVwiLCBmaWVsZDogXCJwYXNzd29yZFwiIH07XG4gICAgcmV0dXJuIG51bGxcbn1cblxuLyoqXG4gKiBIYW5kbGUgc2lnbnVwIHByb2Nlc3NpbmdcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtleHByZXNzLlJlcXVlc3R9IHJlcSBUaGUgcmVxdWVzdCBvYmplY3RcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXNwb25zZX0gcmVzIFRoZSByZXNwb25zZSBvYmplY3RcbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIHNpZ251cChyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICB0cnkge1xuXG4gICAgICAgIGxldCBib2R5ID0gKHJlcS5ib2R5IGFzIG1vZGVscy5BY2NvdW50KVxuICAgICAgICBjb25zdCB2YWxFcnJvciA9IHZhbGlkYXRlQWNjb3VudChib2R5KVxuXG4gICAgICAgIC8vIHNlbmQgdmFsaWRhdGlvbiBlcnJvclxuICAgICAgICBpZiAodmFsRXJyb3IgIT0gbnVsbCkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcihyZXEsIHZhbEVycm9yLm1zZywgdmFsRXJyb3IuZmllbGQsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KCdiYWNrJylcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIGNoZWNrIGlmIGFuIGFjY291bnQgd2l0aCBtYXRjaGluZyBlbWFpbCBleGlzdHNcbiAgICAgICAgbGV0IGV4aXN0aW5nQWNjb3VudCA9IGF3YWl0IEFjY291bnQuZ2V0KHsgZW1haWw6IGJvZHkuZW1haWwgfSlcbiAgICAgICAgaWYgKGV4aXN0aW5nQWNjb3VudCkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcihyZXEsIFwiRW1haWwgaGFzIGFscmVhZHkgYmVlbiByZWdpc3RlcmVkXCIsIFwiZW1haWxcIiwgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gc2V0IGlkLCBjbGllbnQgaWQgYW5kIHNlY3JldFxuICAgICAgICBib2R5LmlkID0gdXVpZDQoKVxuICAgICAgICBib2R5LnBhc3N3b3JkID0gYmNyeXB0Lmhhc2hTeW5jKGJvZHkucGFzc3dvcmQsIDEwKVxuICAgICAgICBib2R5LmNvbmZpcm1hdGlvbl9jb2RlID0gcmFuZG9tLmdlbmVyYXRlKDI4KVxuICAgICAgICBsZXQgYWNjb3VudCA9IGF3YWl0IEFjY291bnQuY3JlYXRlKGJvZHkpXG5cbiAgICAgICAgLy8gc2VuZCBjb25maXJtYXRpb24gbWVzc2FnZVxuICAgICAgICBsZXQgcmVzcCA9IGF3YWl0IEVtYWlsU2VydmljZS5zZW5kQWNjb3VudENvbmZpcm1hdGlvbihib2R5LmVtYWlsLCBhY2NvdW50LnRvSlNPTigpKVxuXG4gICAgICAgIHJlcS5mbGFzaChcImFjY291bnRcIiwgYWNjb3VudClcbiAgICAgICAgcmVzLnJlZGlyZWN0KFwiL2NvbmZpcm1cIilcblxuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihlKVxuICAgICAgICByZXR1cm4gcmVzLnNlcnZlckVycm9yKClcbiAgICB9XG59XG5cbi8qKlxuICogTG9naW4gcHJvY2Vzc2luZ1xuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVxdWVzdH0gcmVxIFRoZSByZXF1ZXN0IG9iamVjdFxuICogQHBhcmFtIHtleHByZXNzLlJlc3BvbnNlfSByZXMgVGhlIHJlc3BvbnNlIG9iamVjdFxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gbG9naW4ocmVxOiBhbnksIHJlczogYW55KSB7XG4gICAgdHJ5IHtcbiAgICAgICAgbGV0IGJvZHkgPSByZXEuYm9keVxuXG4gICAgICAgIGlmICghYm9keS5lbWFpbCkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcihyZXEsICdQbGVhc2UgZW50ZXIgeW91ciBlbWFpbCcsICdlbWFpbCcsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KCdiYWNrJylcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICghYm9keS5wYXNzd29yZCkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcihyZXEsICdQbGVhc2UgZW50ZXIgeW91ciBwYXNzd29yZCcsICdwYXNzd29yZCcsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KCdiYWNrJylcbiAgICAgICAgfVxuXG4gICAgICAgIGxldCBhY2NvdW50ID0gYXdhaXQgQWNjb3VudC5nZXQoeyBlbWFpbDogYm9keS5lbWFpbCB9KVxuICAgICAgICBpZiAoIWFjY291bnQpIHtcbiAgICAgICAgICAgIHNldEZsYXNoRXJyb3IocmVxLCAnWW91ciBlbWFpbCBhbmQgcGFzc3dvcmQgY29tYmluYXRpb24gaXMgbm90IHZhbGlkJywgJycsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KCdiYWNrJylcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICghYmNyeXB0LmNvbXBhcmVTeW5jKGJvZHkucGFzc3dvcmQsIGFjY291bnQucGFzc3dvcmQpKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgJ1lvdXIgZW1haWwgYW5kIHBhc3N3b3JkIGNvbWJpbmF0aW9uIGlzIG5vdCB2YWxpZCcsICcnLCBib2R5KVxuICAgICAgICAgICAgcmV0dXJuIHJlcy5yZWRpcmVjdCgnYmFjaycpXG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIWFjY291bnQuY29uZmlybWVkKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ3dhcm4nLCAnd2FybicpXG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgYFlvdXIgYWNjb3VudCBoYXMgbm90IGJlZW4gY29uZmlybWVkLiBQbGVhc2UgY2hlY2sgeW91ciBlbWFpbCBmb3IgdGhlIGNvbmZpcm1hdGlvbiBtZXNzYWdlLiA8YSBocmVmPVwiL2FjY291bnQvcmVzZW5kLWNvZGUvJHthY2NvdW50LmlkfVwiPlJlc2VuZDwvYT5gLCAnJywgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG5cbiAgICAgICAgcmV0dXJuIHJlcy5zZW5kKGFjY291bnQpXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICByZXR1cm4gcmVzLnNlcnZlckVycm9yKClcbiAgICB9XG59XG5cbi8qKlxuICogVmVyaWZ5IGFuIGFjY291bnRcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtleHByZXNzLlJlcXVlc3R9IHJlcSBUaGUgcmVxdWVzdCBvYmplY3RcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXNwb25zZX0gcmVzIFRoZSByZXNwb25zZSBvYmplY3RcbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIHJlc2VuZENvbmZpcm1hdGlvbihyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICB0cnkge1xuXG4gICAgICAgIGxldCB1c2VySUQgPSByZXEucGFyYW0oXCJpZFwiKVxuICAgICAgICBpZiAoIXVzZXJJRCkge1xuICAgICAgICAgICAgcmV0dXJuIHJlcy5zZXJ2ZXJFcnJvcigpXG4gICAgICAgIH1cblxuICAgICAgICAvLyBnZXQgYWNjb3VudCBieSB1c2VyIGlkXG4gICAgICAgIGxldCBhY2NvdW50ID0gYXdhaXQgQWNjb3VudC5nZXQoeyBpZDogdXNlcklEIH0pXG4gICAgICAgIGlmICghYWNjb3VudCkge1xuICAgICAgICAgICAgcmVxLmZsYXNoKCdlcnJvcicsICdTb3JyeSwgVGhpcyBjb25maXJtYXRpb24gbGluayBpcyBpbnZhbGlkJylcbiAgICAgICAgICAgIHJldHVybiByZXMudmlldygnYWNjb3VudC92ZXJpZnknKVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKGFjY291bnQuY29uZmlybWVkKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ21zZycsICdUaGlzIGFjY291bnQgaGFzIGFscmVhZHkgYmVlbiBjb25maXJtZWQnKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3ZlcmlmeScpXG4gICAgICAgIH1cblxuICAgICAgICAvLyBzZW5kIGNvbmZpcm1hdGlvbiBtZXNzYWdlXG4gICAgICAgIGxldCByZXNwID0gYXdhaXQgRW1haWxTZXJ2aWNlLnNlbmRBY2NvdW50Q29uZmlybWF0aW9uKGFjY291bnQuZW1haWwsIGFjY291bnQudG9KU09OKCkpXG4gICAgICAgIHJlcS5mbGFzaCgnbXNnJywgJ0NvbmZpcm1hdGlvbiBtZXNzYWdlIGhhcyBiZWVuIHNlbnQuIFBsZWFzZSBjaGVjayB5b3VyIGVtYWlsJylcbiAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3ZlcmlmeScpXG5cbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn1cblxuLyoqXG4gKiBTZW5kIHJlc2V0IHBhc3N3b3JkIHRva2VuIHRvIGVtYWlsXG4gKiBcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXF1ZXN0fSByZXEgVGhlIHJlcXVlc3Qgb2JqZWN0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVzcG9uc2V9IHJlcyBUaGUgcmVzcG9uc2Ugb2JqZWN0XG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBmb3Jnb3RQYXNzd29yZChyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICB0cnkge1xuICAgICAgICBsZXQgZW1haWwgPSByZXEuYm9keS5lbWFpbFxuXG4gICAgICAgIGlmICghZW1haWwpIHtcbiAgICAgICAgICAgIHNldEZsYXNoRXJyb3IocmVxLCBcIlBsZWFzZSBlbnRlciB5b3VyIGVtYWlsIGFkZHJlc3NcIiwgXCJlbWFpbFwiLCByZXEuYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMudmlldygnYWNjb3VudC9mb3Jnb3RfcGFzc3dvcmQnKVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gZ2V0IGFjY291bnQgYnkgZW1haWxcbiAgICAgICAgbGV0IGFjY291bnQgPSBhd2FpdCBBY2NvdW50LmdldCh7IGVtYWlsIH0pXG4gICAgICAgIGlmICghYWNjb3VudCkge1xuICAgICAgICAgICAgcmVxLmZsYXNoKCdlcnJvcicsIGBObyBhY2NvdW50IG1hdGNoaW5nIHRoZSBlbWFpbCAnJHtyZXEuYm9keS5lbWFpbH0nYClcbiAgICAgICAgICAgIHJldHVybiByZXMudmlldygnYWNjb3VudC9mb3Jnb3RfcGFzc3dvcmQnKVxuICAgICAgICB9XG5cbiAgICAgICAgbGV0IHJlc2V0VG9rZW4gPSBhd2FpdCBSZXNldFRva2VuLmNyZWF0ZSh7IHRva2VuOiBzaGExKHV1aWQ0KCkpLCBhY2NvdW50OiBhY2NvdW50LmlkIH0pXG5cbiAgICAgICAgLy8gc2VuZCBwYXNzd29yZCByZXNldCBtZXNzYWdlXG4gICAgICAgIGFjY291bnQgPSBhY2NvdW50LnRvSlNPTigpXG4gICAgICAgIGFjY291bnQudG9rZW4gPSByZXNldFRva2VuLnRva2VuXG4gICAgICAgIGxldCByZXNwID0gYXdhaXQgRW1haWxTZXJ2aWNlLnNlbmRQYXNzd29yZFJlc2V0RW1haWwoYWNjb3VudC5lbWFpbCwgYWNjb3VudClcbiAgICAgICAgcmVxLmZsYXNoKCdtc2cnLCAnQSByZXNldCBtZXNzYWdlIGhhcyBiZWVuIHNlbnQgdG8geW91JylcbiAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L2ZvcmdvdF9wYXNzd29yZCcpXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICByZXR1cm4gcmVzLnNlcnZlckVycm9yKClcbiAgICB9XG59XG5cbi8qKlxuICogUmVzZXQgYSBwYXNzd29yZCBwYWdlXG4gKiBcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXF1ZXN0fSByZXEgVGhlIHJlcXVlc3Qgb2JqZWN0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVzcG9uc2V9IHJlcyBUaGUgcmVzcG9uc2Ugb2JqZWN0XG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiByZXNldFBhc3N3b3JkUGFnZShyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICB0cnkge1xuICAgICAgICBsZXQgcmVzZXRUb2tlbiA9IGF3YWl0IFJlc2V0VG9rZW4uZ2V0KHsgdG9rZW46IHJlcS5wYXJhbShcInRva2VuXCIpIH0pXG4gICAgICAgIGlmICghcmVzZXRUb2tlbikge1xuICAgICAgICAgICAgcmVxLmZsYXNoKCdlcnJvcjInLCBgUmVzZXQgbGluayBpcyBpbnZhbGlkIG9yIGhhcyBleHBpcmVkYClcbiAgICAgICAgICAgIHJldHVybiByZXMudmlldygnYWNjb3VudC9yZXNldF9wYXNzd29yZCcpXG4gICAgICAgIH1cblxuICAgICAgICBpZiAocmVzZXRUb2tlbi51c2VkKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yMicsIGBSZXNldCBsaW5rIGlzIGludmFsaWQgb3IgaGFzIGV4cGlyZWRgKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3Jlc2V0X3Bhc3N3b3JkJylcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICghbW9tZW50KHJlc2V0VG9rZW4uY3JlYXRlZF9hdCkuYWRkKDIsIFwiZGF5XCIpLnV0YygpLmlzQWZ0ZXIobW9tZW50KCkpKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yMicsIGBSZXNldCBsaW5rIGlzIGludmFsaWQgb3IgaGFzIGV4cGlyZWRgKVxuICAgICAgICB9XG5cbiAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3Jlc2V0X3Bhc3N3b3JkJylcblxuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgcmV0dXJuIHJlcy5zZXJ2ZXJFcnJvcigpXG4gICAgfVxufVxuXG4vKipcbiAqIFJlc2V0IGEgcGFzc3dvcmRcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtleHByZXNzLlJlcXVlc3R9IHJlcSBUaGUgcmVxdWVzdCBvYmplY3RcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXNwb25zZX0gcmVzIFRoZSByZXNwb25zZSBvYmplY3RcbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIHJlc2V0UGFzc3dvcmQocmVxOiBhbnksIHJlczogYW55KSB7XG4gICAgdHJ5IHtcblxuICAgICAgICBsZXQgcmVzZXRUb2tlbiA9IGF3YWl0IFJlc2V0VG9rZW4uZ2V0KHsgdG9rZW46IHJlcS5wYXJhbShcInRva2VuXCIpIH0pXG4gICAgICAgIGlmICghcmVzZXRUb2tlbikge1xuICAgICAgICAgICAgcmVxLmZsYXNoKCdlcnJvcjInLCBgUmVzZXQgbGluayBpcyBpbnZhbGlkIG9yIGhhcyBleHBpcmVkYClcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QocmVxLm9yaWdpbmFsVXJsKVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHJlc2V0VG9rZW4udXNlZCkge1xuICAgICAgICAgICAgcmVxLmZsYXNoKCdlcnJvcjInLCBgUmVzZXQgbGluayBpcyBpbnZhbGlkIG9yIGhhcyBleHBpcmVkYClcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QocmVxLm9yaWdpbmFsVXJsKVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKCFtb21lbnQocmVzZXRUb2tlbi5jcmVhdGVkX2F0KS5hZGQoMiwgXCJkYXlcIikudXRjKCkuaXNBZnRlcihtb21lbnQoKSkpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnZXJyb3IyJywgYFJlc2V0IGxpbmsgaXMgaW52YWxpZCBvciBoYXMgZXhwaXJlZGApXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KHJlcS5vcmlnaW5hbFVybClcbiAgICAgICAgfVxuXG4gICAgICAgIGxldCBib2R5ID0gcmVxLmJvZHlcblxuICAgICAgICAvLyB2YWxpZGF0aW9uXG4gICAgICAgIGlmICh2YWxpZGF0b3IuaXNFbXB0eShib2R5LnBhc3N3b3JkKSkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcihyZXEsIGBQYXNzd29yZCBpcyByZXF1aXJlZGAsICdwYXNzd29yZCcsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KHJlcS5vcmlnaW5hbFVybClcbiAgICAgICAgfSBlbHNlIGlmICghdmFsaWRhdG9yLmlzTGVuZ3RoKGJvZHkucGFzc3dvcmQsIHsgbWluOiA4LCBtYXg6IDEyMCB9KSkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcihyZXEsIGBQbGVhc2UgZW50ZXIgYSBwYXNzd29yZCBiZXR3ZWVuIDggYW5kIDEyMCBjaGFyYWN0ZXJzLCB0aGUgdHJpY2tpZXIgdGhlIGJldHRlciFgLCAncGFzc3dvcmQnLCBib2R5KVxuICAgICAgICAgICAgcmV0dXJuIHJlcy5yZWRpcmVjdChyZXEub3JpZ2luYWxVcmwpXG4gICAgICAgIH1cblxuICAgICAgICBsZXQgYWNjb3VudCA9IGF3YWl0IEFjY291bnQuZ2V0KHsgaWQ6IHJlc2V0VG9rZW4uYWNjb3VudCB9KVxuICAgICAgICBpZiAoIWFjY291bnQpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnZXJyb3IyJywgYFJlc2V0IGxpbmsgaXMgaW52YWxpZCBvciBoYXMgZXhwaXJlZCdgKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy5yZWRpcmVjdChyZXEub3JpZ2luYWxVcmwpXG4gICAgICAgIH1cblxuICAgICAgICAvLyBlbnN1cmUgbmV3IHBhc3N3b3JkIGlzIG5vdCB0aGUgc2FtZSBhcyB0aGUgb2xkIG9uZVxuICAgICAgICBpZiAoYmNyeXB0LmNvbXBhcmVTeW5jKGJvZHkucGFzc3dvcmQsIGFjY291bnQucGFzc3dvcmQpKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgYE5ldyBwYXNzd29yZCBjYW5ub3QgYmUgc2ltaWxhciB0byBhIHByZXZpb3VzbHkgdXNlZCBwYXNzd29yZGAsICdwYXNzd29yZCcsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KHJlcS5vcmlnaW5hbFVybClcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIHNldCByZXNldCB0b2tlbiBhcyB1c2VkXG4gICAgICAgIHJlc2V0VG9rZW4uc2V0KCd1c2VkJywgdHJ1ZSlcbiAgICAgICAgYXdhaXQgcmVzZXRUb2tlbi5zYXZlKClcblxuICAgICAgICAvLyB1cGRhdGUgcGFzc3dvcmRcbiAgICAgICAgbGV0IG5ld1Bhc3N3b3JkSGFzaCA9IGJjcnlwdC5oYXNoU3luYyhib2R5LnBhc3N3b3JkLCAxMClcbiAgICAgICAgYWNjb3VudC5zZXQoJ3Bhc3N3b3JkJywgbmV3UGFzc3dvcmRIYXNoKVxuICAgICAgICBhd2FpdCBhY2NvdW50LnNhdmUoKVxuXG4gICAgICAgIHJldHVybiByZXMudmlldygnYWNjb3VudC9yZXNldF9wYXNzd29yZCcsIHsgbXNnOiAnUGFzc3dvcmQgc3VjY2Vzc2Z1bGx5IGNoYW5nZWQnIH0pXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICByZXR1cm4gcmVzLnNlcnZlckVycm9yKClcbiAgICB9XG59XG5cbi8qKlxuICogVmVyaWZ5IGFuIGFjY291bnRcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtleHByZXNzLlJlcXVlc3R9IHJlcSBUaGUgcmVxdWVzdCBvYmplY3RcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXNwb25zZX0gcmVzIFRoZSByZXNwb25zZSBvYmplY3RcbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIHZlcmlmeShyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICB0cnkge1xuXG4gICAgICAgIGxldCBjb2RlID0gcmVxLnBhcmFtKFwiY29kZVwiKVxuICAgICAgICBpZiAoIWNvZGUpIHtcbiAgICAgICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gZ2V0IGFjY291bnQgYnkgY29uZmlybWF0aW9uIGNvZGVcbiAgICAgICAgbGV0IGFjY291bnQgPSBhd2FpdCBBY2NvdW50LmdldCh7IGNvbmZpcm1hdGlvbl9jb2RlOiBjb2RlIH0pXG4gICAgICAgIGlmICghYWNjb3VudCkge1xuICAgICAgICAgICAgcmVxLmZsYXNoKCdlcnJvcicsICdTb3JyeSwgVGhpcyBjb25maXJtYXRpb24gbGluayBpcyBpbnZhbGlkJylcbiAgICAgICAgICAgIHJldHVybiByZXMudmlldygnYWNjb3VudC92ZXJpZnknKVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKGFjY291bnQuY29uZmlybWVkKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ21zZycsICdUaGlzIGFjY291bnQgaGFzIGFscmVhZHkgYmVlbiBjb25maXJtZWQnKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3ZlcmlmeScpXG4gICAgICAgIH1cblxuICAgICAgICBhY2NvdW50LnNldCgnY29uZmlybWVkJywgdHJ1ZSlcbiAgICAgICAgYXdhaXQgYWNjb3VudC5zYXZlKClcbiAgICAgICAgcmVxLmZsYXNoKCdtc2cnLCAnR3JlYXQhIFlvdXIgYWNjb3VudCBoYXMgYmVlbiBjb25maXJtZWQnKVxuICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvdmVyaWZ5JylcblxuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihlKVxuICAgICAgICByZXR1cm4gcmVzLnNlcnZlckVycm9yKClcbiAgICB9XG59Il19
