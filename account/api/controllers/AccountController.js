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
const Random = require("random-js");
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
            body.client_id = random.generate(32);
            body.client_secret = random.generate(Random.integer(27, 32)(Random.engines.mt19937().autoSeed()));
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
            let resp = yield EmailService.sendPasswordResetEmail('kennedyidialu@gmail.com', account);
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

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9jb250cm9sbGVycy9BY2NvdW50Q29udHJvbGxlci50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7O0FBR0EsdUNBQXdDO0FBQ3hDLDZDQUE4QztBQUM5QyxtREFBb0Q7QUFDcEQsK0JBQWdDO0FBQ2hDLHVDQUF3QztBQUN4QyxvQ0FBcUM7QUFDckMsaUNBQWtDO0FBQ2xDLDZCQUE4QjtBQUM5QixpQ0FBa0M7QUFFbEM7Ozs7Ozs7R0FPRztBQUNILHVCQUF1QixHQUFRLEVBQUUsR0FBVyxFQUFFLEtBQWEsRUFBRSxJQUFTO0lBQ2xFLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxFQUFFLEdBQUcsQ0FBQyxDQUFBO0lBQ3ZCLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxFQUFFLEtBQUssQ0FBQyxDQUFBO0lBQ3pCLEdBQUcsQ0FBQyxLQUFLLENBQUMsTUFBTSxFQUFFLElBQUksQ0FBQyxDQUFBO0FBQzNCLENBQUM7QUFHRDs7Ozs7R0FLRztBQUNILHlCQUFnQyxPQUFnQjtJQUM1QyxFQUFFLENBQUMsQ0FBQyxTQUFTLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxVQUFVLENBQUMsQ0FBQztRQUFDLE1BQU0sQ0FBQyxFQUFFLEdBQUcsRUFBRSx3QkFBd0IsRUFBRSxLQUFLLEVBQUUsWUFBWSxFQUFFLENBQUM7SUFDekcsRUFBRSxDQUFDLENBQUMsU0FBUyxDQUFDLE9BQU8sQ0FBQyxPQUFPLENBQUMsU0FBUyxDQUFDLENBQUM7UUFBQyxNQUFNLENBQUMsRUFBRSxHQUFHLEVBQUUsdUJBQXVCLEVBQUUsS0FBSyxFQUFFLFdBQVcsRUFBRSxDQUFDO0lBQ3RHLEVBQUUsQ0FBQyxDQUFDLFNBQVMsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLEtBQUssQ0FBQyxDQUFDO1FBQUMsTUFBTSxDQUFDLEVBQUUsR0FBRyxFQUFFLG1CQUFtQixFQUFFLEtBQUssRUFBRSxPQUFPLEVBQUUsQ0FBQztJQUMxRixFQUFFLENBQUMsQ0FBQyxDQUFDLFNB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SUFDTCxDQUFDO0NBQUE7QUEzQkQsd0NBMkJDO0FBRUQ7Ozs7OztHQU1HO0FBQ0gsMkJBQXdDLEdBQVEsRUFBRSxHQUFROztRQUN0RCxJQUFJLENBQUM7WUFDRCxJQUFJLFVBQVUsR0FBRyxNQUFNLFVBQVUsQ0FBQyxHQUFHLENBQUMsRUFBRSxLQUFLLEVBQUUsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsRUFBRSxDQUFDLENBQUE7WUFDcEUsRUFBRSxDQUFDLENBQUMsQ0FBQyxVQUFVLENBQUMsQ0FBQyxDQUFDO2dCQUNkLEdBQUcsQ0FBQyxLQUFLLENBQUMsUUFBUSxFQUFFLHNDQUFzQyxDQUFDLENBQUE7Z0JBQzNELE1BQU0sQ0FBQyxHQUFHLENBQUMsSUFBSSxDQUFDLHdCQUF3QixDQUFDLENBQUE7WUFDN0MsQ0FBQztZQUVELEVBQUUsQ0FBQyxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDO2dCQUNsQixHQUFHLENBQUMsS0FBSyxDQUFDLFFBQVEsRUFBRSxzQ0FBc0MsQ0FBQyxDQUFBO2dCQUMzRCxNQUFNLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyx3QkFBd0IsQ0FBQyxDQUFBO1lBQzdDLENBQUM7WUFFRCxFQUFFLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxVQUFVLENBQUMsVUFBVSxDQUFDLENBQUMsR0FBRyxDQUFDLENBQUMsRUFBRSxLQUFLLENBQUMsQ0FBQyxHQUFHLEVBQUUsQ0FBQyxPQUFPLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQyxDQUFDLENBQUM7Z0JBQ3ZFLEdBQUcsQ0FBQyxLQUFLLENBQUMsUUFBUSxFQUFFLHNDQUFzQyxDQUFDLENBQUE7WUFDL0QsQ0FBQztZQUVELE1BQU0sQ0FBQyxHQUFHLENBQUMsSUFBSSxDQUFDLHdCQUF3QixDQUFDLENBQUE7UUFFN0MsQ0FBQztRQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUM7WUFDVCxNQUFNLENBQUMsR0FBRyxDQUFDLFdBQVcsRUFBRSxDQUFBO1FBQzVCLENBQUM7SUFDTCxD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QztnQkFDWCxHQUFHLENBQUMsS0FBSyxDQUFDLE9BQU8sRUFBRSwwQ0FBMEMsQ0FBQyxDQUFBO2dCQUM5RCxNQUFNLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxnQkFBZ0IsQ0FBQyxDQUFBO1lBQ3JDLENBQUM7WUFFRCxFQUFFLENBQUMsQ0FBQyxPQUFPLENBQUMsU0FBUyxDQUFDLENBQUMsQ0FBQztnQkFDcEIsR0FBRyxDQUFDLEtBQUssQ0FBQyxLQUFLLEVBQUUseUNBQXlDLENBQUMsQ0FBQTtnQkFDM0QsTUFBTSxDQUFDLEdBQUcsQ0FBQyxJQUFJLENBQUMsZ0JBQWdCLENBQUMsQ0FBQTtZQUNyQyxDQUFDO1lBRUQsT0FBTyxDQUFDLEdBQUcsQ0FBQyxXQUFXLEVBQUUsSUFBSSxDQUFDLENBQUE7WUFDOUIsTUFBTSxPQUFPLENBQUMsSUFBSSxFQUFFLENBQUE7WUFDcEIsR0FBRyxDQUFDLEtBQUssQ0FBQyxLQUFLLEVBQUUsd0NBQXdDLENBQUMsQ0FBQTtZQUMxRCxNQUFNLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxnQkFBZ0IsQ0FBQyxDQUFBO1FBRXJDLENBQUM7UUFBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDO1lBQ1QsT0FBTyxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUNoQixNQUFNLENBQUMsR0FBRyxDQUFDLFdBQVcsRUFBRSxDQUFBO1FBQzVCLENBQUM7SUFDTCxDQUFDO0NBQUE7QUE3QkQsd0JBNkJDIiwiZmlsZSI6ImFwaS9jb250cm9sbGVycy9BY2NvdW50Q29udHJvbGxlci5qcyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IHNldEZsYWdzRnJvbVN0cmluZyB9IGZyb20gJ3Y4JztcbmltcG9ydCB7IGFjY2VzcyB9IGZyb20gJ2ZzJztcbmltcG9ydCBleHByZXNzID0gcmVxdWlyZSgnZXhwcmVzcycpO1xuaW1wb3J0IHZhbGlkYXRvciA9IHJlcXVpcmUoJ3ZhbGlkYXRvcicpO1xuaW1wb3J0IEFjY291bnQgPSByZXF1aXJlKCcuLi9tb2RlbHMvQWNjb3VudCcpO1xuaW1wb3J0IFJlc2V0VG9rZW4gPSByZXF1aXJlKCcuLi9tb2RlbHMvUmVzZXRUb2tlbicpO1xuaW1wb3J0IHV1aWQ0ID0gcmVxdWlyZSgndXVpZDQnKTtcbmltcG9ydCByYW5kb20gPSByZXF1aXJlKCdyYW5kb21zdHJpbmcnKTtcbmltcG9ydCBSYW5kb20gPSByZXF1aXJlKCdyYW5kb20tanMnKTtcbmltcG9ydCBiY3J5cHQgPSByZXF1aXJlKCdiY3J5cHQnKTtcbmltcG9ydCBzaGExID0gcmVxdWlyZSgnc2hhMScpO1xuaW1wb3J0IG1vbWVudCA9IHJlcXVpcmUoJ21vbWVudCcpO1xuXG4vKipcbiAqIFNldHMgZmxhc2ggZXJyb3IgbWVzc2FnZSBmb3IgYSBmaWVsZFxuICogXG4gKiBAcGFyYW0geyp9IHJlcSBSZXF1ZXN0IG9iamVjdFxuICogQHBhcmFtIHtzdHJpbmd9IG1zZyBUaGUgZXJyb3IgbWVzc2FnZVxuICogQHBhcmFtIHtzdHJpbmd9IGZpZWxkIFRoZSBmaWVsZCBuYW1lXG4gKiBAcGFyYW0geyp9IGJvZHkgUmVzcG9uc2UgYm9keVxuICovXG5mdW5jdGlvbiBzZXRGbGFzaEVycm9yKHJlcTogYW55LCBtc2c6IHN0cmluZywgZmllbGQ6IHN0cmluZywgYm9keTogYW55KSB7XG4gICAgcmVxLmZsYXNoKFwiZXJyb3JcIiwgbXNnKVxuICAgIHJlcS5mbGFzaChcImZpZWxkXCIsIGZpZWxkKVxuICAgIHJlcS5mbGFzaChcImZvcm1cIiwgYm9keSlcbn1cblxuXG4vKipcbiAqIFZhbGlkYXRlcyByZXMuYm9keSBjb250YWluaW5nIGFjY291bnQgZGV0YWlsc1xuICogQGV4cG9ydFxuICogQHBhcmFtIHtBY2NvdW50fSBhY2NvdW50IFxuICogQHJldHVybnMge1ZhbGlkYXRpb25FcnJvcn0gXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiB2YWxpZGF0ZUFjY291bnQoYWNjb3VudDogQWNjb3VudCk6IFZhbGlkYXRpb25FcnJvciB7XG4gICAgaWYgKHZhbGlkYXRvci5pc0VtcHR5KGFjY291bnQuZmlyc3RfbmFtZSkpIHJldHVybiB7IG1zZzogXCJGaXJzdCBuYW1lIGlzIHJlcXVpcmVkXCIsIGZpZWxkOiBcImZpcnN0X25hbWVcIiB9O1xuICAgIGlmICh2YWxpZGF0b3IuaXNFbXB0eShhY2NvdW50Lmxhc3RfbmFtZSkpIHJldHVybiB7IG1zZzogXCJMYXN0IG5hbWUgaXMgcmVxdWlyZWRcIiwgZmllbGQ6IFwibGFzdF9uYW1lXCIgfTtcbiAgICBpZiAodmFsaWRhdG9yLmlzRW1wdHkoYWNjb3VudC5lbWFpbCkpIHJldHVybiB7IG1zZzogXCJFbWFpbCBpcyByZXF1aXJlZFwiLCBmaWVsZDogXCJlbWFpbFwiIH07XG4gICAgaWYgKCF2YWxpZGF0b3IuaXNFbWFpbChhY2NvdW50LmVtYWlsKSkgcmV0dXJuIHsgbXNnOiBcIkVtYWlsIGRvZXMgbm90IGFwcGVhciB0byBiZSB2YWxpZFwiLCBmaWVsZDogXCJlbWFpbFwiIH07XG4gICAgaWYgKHZhbGlkYXRvci5pc0VtcHR5KGFjY291bnQucGFzc3dvcmQpKSByZXR1cm4geyBtc2c6IFwiUGFzc3dvcmQgaXMgcmVxdWlyZWRcIiwgZmllbGQ6IFwicGFzc3dvcmRcIiB9O1xuICAgIGlmICghdmFsaWRhdG9yLmlzTGVuZ3RoKGFjY291bnQucGFzc3dvcmQsIHsgbWluOiA4LCBtYXg6IDEyMCB9KSkgcmV0dXJuIHsgbXNnOiBcIlBsZWFzZSBlbnRlciBhIHBhc3N3b3JkIGJldHdlZW4gOCBhbmQgMTIwIGNoYXJhY3RlcnMsIHRoZSB0cmlja2llciB0aGUgYmV0dGVyIVwiLCBmaWVsZDogXCJwYXNzd29yZFwiIH07XG4gICAgcmV0dXJuIG51bGxcbn1cblxuLyoqXG4gKiBIYW5kbGUgc2lnbnVwIHByb2Nlc3NpbmdcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtleHByZXNzLlJlcXVlc3R9IHJlcSBUaGUgcmVxdWVzdCBvYmplY3RcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXNwb25zZX0gcmVzIFRoZSByZXNwb25zZSBvYmplY3RcbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIHNpZ251cChyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICB0cnkge1xuXG4gICAgICAgIGxldCBib2R5ID0gKHJlcS5ib2R5IGFzIEFjY291bnQpXG4gICAgICAgIGNvbnN0IHZhbEVycm9yID0gdmFsaWRhdGVBY2NvdW50KGJvZHkpXG5cbiAgICAgICAgLy8gc2VuZCB2YWxpZGF0aW9uIGVycm9yXG4gICAgICAgIGlmICh2YWxFcnJvciAhPSBudWxsKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgdmFsRXJyb3IubXNnLCB2YWxFcnJvci5maWVsZCwgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gY2hlY2sgaWYgYW4gYWNjb3VudCB3aXRoIG1hdGNoaW5nIGVtYWlsIGV4aXN0c1xuICAgICAgICBsZXQgZXhpc3RpbmdBY2NvdW50ID0gYXdhaXQgQWNjb3VudC5nZXQoeyBlbWFpbDogYm9keS5lbWFpbCB9KVxuICAgICAgICBpZiAoZXhpc3RpbmdBY2NvdW50KSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgXCJFbWFpbCBoYXMgYWxyZWFkeSBiZWVuIHJlZ2lzdGVyZWRcIiwgXCJlbWFpbFwiLCBib2R5KVxuICAgICAgICAgICAgcmV0dXJuIHJlcy5yZWRpcmVjdCgnYmFjaycpXG4gICAgICAgIH1cblxuICAgICAgICAvLyBzZXQgaWQsIGNsaWVudCBpZCBhbmQgc2VjcmV0XG4gICAgICAgIGJvZHkuaWQgPSB1dWlkNCgpXG4gICAgICAgIGJvZHkuY2xpZW50X2lkID0gcmFuZG9tLmdlbmVyYXRlKDMyKVxuICAgICAgICBib2R5LmNsaWVudF9zZWNyZXQgPSByYW5kb20uZ2VuZXJhdGUoUmFuZG9tLmludGVnZXIoMjcsIDMyKShSYW5kb20uZW5naW5lcy5tdDE5OTM3KCkuYXV0b1NlZWQoKSkpXG4gICAgICAgIGJvZHkucGFzc3dvcmQgPSBiY3J5cHQuaGFzaFN5bmMoYm9keS5wYXNzd29yZCwgMTApXG4gICAgICAgIGJvZHkuY29uZmlybWF0aW9uX2NvZGUgPSByYW5kb20uZ2VuZXJhdGUoMjgpXG4gICAgICAgIGxldCBhY2NvdW50ID0gYXdhaXQgQWNjb3VudC5jcmVhdGUoYm9keSlcblxuICAgICAgICAvLyBzZW5kIGNvbmZpcm1hdGlvbiBtZXNzYWdlXG4gICAgICAgIGxldCByZXNwID0gYXdhaXQgRW1haWxTZXJ2aWNlLnNlbmRBY2NvdW50Q29uZmlybWF0aW9uKGJvZHkuZW1haWwsIGFjY291bnQudG9KU09OKCkpXG5cbiAgICAgICAgcmVxLmZsYXNoKFwiYWNjb3VudFwiLCBhY2NvdW50KVxuICAgICAgICByZXMucmVkaXJlY3QoXCIvY29uZmlybVwiKVxuXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zb2xlLmVycm9yKGUpXG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn1cblxuLyoqXG4gKiBMb2dpbiBwcm9jZXNzaW5nXG4gKiBcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXF1ZXN0fSByZXEgVGhlIHJlcXVlc3Qgb2JqZWN0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVzcG9uc2V9IHJlcyBUaGUgcmVzcG9uc2Ugb2JqZWN0XG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBsb2dpbihyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICB0cnkge1xuICAgICAgICBsZXQgYm9keSA9IHJlcS5ib2R5XG5cbiAgICAgICAgaWYgKCFib2R5LmVtYWlsKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgJ1BsZWFzZSBlbnRlciB5b3VyIGVtYWlsJywgJ2VtYWlsJywgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKCFib2R5LnBhc3N3b3JkKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgJ1BsZWFzZSBlbnRlciB5b3VyIHBhc3N3b3JkJywgJ3Bhc3N3b3JkJywgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG5cbiAgICAgICAgbGV0IGFjY291bnQgPSBhd2FpdCBBY2NvdW50LmdldCh7IGVtYWlsOiBib2R5LmVtYWlsIH0pXG4gICAgICAgIGlmICghYWNjb3VudCkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcihyZXEsICdZb3VyIGVtYWlsIGFuZCBwYXNzd29yZCBjb21iaW5hdGlvbiBpcyBub3QgdmFsaWQnLCAnJywgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKCFiY3J5cHQuY29tcGFyZVN5bmMoYm9keS5wYXNzd29yZCwgYWNjb3VudC5wYXNzd29yZCkpIHtcbiAgICAgICAgICAgIHNldEZsYXNoRXJyb3IocmVxLCAnWW91ciBlbWFpbCBhbmQgcGFzc3dvcmQgY29tYmluYXRpb24gaXMgbm90IHZhbGlkJywgJycsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KCdiYWNrJylcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICghYWNjb3VudC5jb25maXJtZWQpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnd2FybicsICd3YXJuJylcbiAgICAgICAgICAgIHNldEZsYXNoRXJyb3IocmVxLCBgWW91ciBhY2NvdW50IGhhcyBub3QgYmVlbiBjb25maXJtZWQuIFBsZWFzZSBjaGVjayB5b3VyIGVtYWlsIGZvciB0aGUgY29uZmlybWF0aW9uIG1lc3NhZ2UuIDxhIGhyZWY9XCIvYWNjb3VudC9yZXNlbmQtY29kZS8ke2FjY291bnQuaWR9XCI+UmVzZW5kPC9hPmAsICcnLCBib2R5KVxuICAgICAgICAgICAgcmV0dXJuIHJlcy5yZWRpcmVjdCgnYmFjaycpXG4gICAgICAgIH1cblxuICAgICAgICByZXR1cm4gcmVzLnNlbmQoYWNjb3VudClcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn1cblxuLyoqXG4gKiBWZXJpZnkgYW4gYWNjb3VudFxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVxdWVzdH0gcmVxIFRoZSByZXF1ZXN0IG9iamVjdFxuICogQHBhcmFtIHtleHByZXNzLlJlc3BvbnNlfSByZXMgVGhlIHJlc3BvbnNlIG9iamVjdFxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gcmVzZW5kQ29uZmlybWF0aW9uKHJlcTogYW55LCByZXM6IGFueSkge1xuICAgIHRyeSB7XG5cbiAgICAgICAgbGV0IHVzZXJJRCA9IHJlcS5wYXJhbShcImlkXCIpXG4gICAgICAgIGlmICghdXNlcklEKSB7XG4gICAgICAgICAgICByZXR1cm4gcmVzLnNlcnZlckVycm9yKClcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIGdldCBhY2NvdW50IGJ5IHVzZXIgaWRcbiAgICAgICAgbGV0IGFjY291bnQgPSBhd2FpdCBBY2NvdW50LmdldCh7IGlkOiB1c2VySUQgfSlcbiAgICAgICAgaWYgKCFhY2NvdW50KSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yJywgJ1NvcnJ5LCBUaGlzIGNvbmZpcm1hdGlvbiBsaW5rIGlzIGludmFsaWQnKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3ZlcmlmeScpXG4gICAgICAgIH1cblxuICAgICAgICBpZiAoYWNjb3VudC5jb25maXJtZWQpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnbXNnJywgJ1RoaXMgYWNjb3VudCBoYXMgYWxyZWFkeSBiZWVuIGNvbmZpcm1lZCcpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvdmVyaWZ5JylcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIHNlbmQgY29uZmlybWF0aW9uIG1lc3NhZ2VcbiAgICAgICAgbGV0IHJlc3AgPSBhd2FpdCBFbWFpbFNlcnZpY2Uuc2VuZEFjY291bnRDb25maXJtYXRpb24oYWNjb3VudC5lbWFpbCwgYWNjb3VudC50b0pTT04oKSlcbiAgICAgICAgcmVxLmZsYXNoKCdtc2cnLCAnQ29uZmlybWF0aW9uIG1lc3NhZ2UgaGFzIGJlZW4gc2VudC4gUGxlYXNlIGNoZWNrIHlvdXIgZW1haWwnKVxuICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvdmVyaWZ5JylcblxuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgcmV0dXJuIHJlcy5zZXJ2ZXJFcnJvcigpXG4gICAgfVxufVxuXG4vKipcbiAqIFNlbmQgcmVzZXQgcGFzc3dvcmQgdG9rZW4gdG8gZW1haWxcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtleHByZXNzLlJlcXVlc3R9IHJlcSBUaGUgcmVxdWVzdCBvYmplY3RcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXNwb25zZX0gcmVzIFRoZSByZXNwb25zZSBvYmplY3RcbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIGZvcmdvdFBhc3N3b3JkKHJlcTogYW55LCByZXM6IGFueSkge1xuICAgIHRyeSB7XG4gICAgICAgIGxldCBlbWFpbCA9IHJlcS5ib2R5LmVtYWlsXG5cbiAgICAgICAgaWYgKCFlbWFpbCkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcihyZXEsIFwiUGxlYXNlIGVudGVyIHlvdXIgZW1haWwgYWRkcmVzc1wiLCBcImVtYWlsXCIsIHJlcS5ib2R5KVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L2ZvcmdvdF9wYXNzd29yZCcpXG4gICAgICAgIH1cblxuICAgICAgICAvLyBnZXQgYWNjb3VudCBieSBlbWFpbFxuICAgICAgICBsZXQgYWNjb3VudCA9IGF3YWl0IEFjY291bnQuZ2V0KHsgZW1haWwgfSlcbiAgICAgICAgaWYgKCFhY2NvdW50KSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yJywgYE5vIGFjY291bnQgbWF0Y2hpbmcgdGhlIGVtYWlsICcke3JlcS5ib2R5LmVtYWlsfSdgKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L2ZvcmdvdF9wYXNzd29yZCcpXG4gICAgICAgIH1cblxuICAgICAgICBsZXQgcmVzZXRUb2tlbiA9IGF3YWl0IFJlc2V0VG9rZW4uY3JlYXRlKHsgdG9rZW46IHNoYTEodXVpZDQoKSksIGFjY291bnQ6IGFjY291bnQuaWQgfSlcblxuICAgICAgICAvLyBzZW5kIHBhc3N3b3JkIHJlc2V0IG1lc3NhZ2VcbiAgICAgICAgYWNjb3VudCA9IGFjY291bnQudG9KU09OKClcbiAgICAgICAgYWNjb3VudC50b2tlbiA9IHJlc2V0VG9rZW4udG9rZW5cbiAgICAgICAgbGV0IHJlc3AgPSBhd2FpdCBFbWFpbFNlcnZpY2Uuc2VuZFBhc3N3b3JkUmVzZXRFbWFpbCgna2VubmVkeWlkaWFsdUBnbWFpbC5jb20nLCBhY2NvdW50KVxuICAgICAgICByZXEuZmxhc2goJ21zZycsICdBIHJlc2V0IG1lc3NhZ2UgaGFzIGJlZW4gc2VudCB0byB5b3UnKVxuICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvZm9yZ290X3Bhc3N3b3JkJylcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn1cblxuLyoqXG4gKiBSZXNldCBhIHBhc3N3b3JkIHBhZ2VcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtleHByZXNzLlJlcXVlc3R9IHJlcSBUaGUgcmVxdWVzdCBvYmplY3RcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXNwb25zZX0gcmVzIFRoZSByZXNwb25zZSBvYmplY3RcbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIHJlc2V0UGFzc3dvcmRQYWdlKHJlcTogYW55LCByZXM6IGFueSkge1xuICAgIHRyeSB7XG4gICAgICAgIGxldCByZXNldFRva2VuID0gYXdhaXQgUmVzZXRUb2tlbi5nZXQoeyB0b2tlbjogcmVxLnBhcmFtKFwidG9rZW5cIikgfSlcbiAgICAgICAgaWYgKCFyZXNldFRva2VuKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yMicsIGBSZXNldCBsaW5rIGlzIGludmFsaWQgb3IgaGFzIGV4cGlyZWRgKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3Jlc2V0X3Bhc3N3b3JkJylcbiAgICAgICAgfVxuXG4gICAgICAgIGlmIChyZXNldFRva2VuLnVzZWQpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnZXJyb3IyJywgYFJlc2V0IGxpbmsgaXMgaW52YWxpZCBvciBoYXMgZXhwaXJlZGApXG4gICAgICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvcmVzZXRfcGFzc3dvcmQnKVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKCFtb21lbnQocmVzZXRUb2tlbi5jcmVhdGVkX2F0KS5hZGQoMiwgXCJkYXlcIikudXRjKCkuaXNBZnRlcihtb21lbnQoKSkpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnZXJyb3IyJywgYFJlc2V0IGxpbmsgaXMgaW52YWxpZCBvciBoYXMgZXhwaXJlZGApXG4gICAgICAgIH1cblxuICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvcmVzZXRfcGFzc3dvcmQnKVxuXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICByZXR1cm4gcmVzLnNlcnZlckVycm9yKClcbiAgICB9XG59XG5cbi8qKlxuICogUmVzZXQgYSBwYXNzd29yZFxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVxdWVzdH0gcmVxIFRoZSByZXF1ZXN0IG9iamVjdFxuICogQHBhcmFtIHtleHByZXNzLlJlc3BvbnNlfSByZXMgVGhlIHJlc3BvbnNlIG9iamVjdFxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gcmVzZXRQYXNzd29yZChyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICB0cnkge1xuXG4gICAgICAgIGxldCByZXNldFRva2VuID0gYXdhaXQgUmVzZXRUb2tlbi5nZXQoeyB0b2tlbjogcmVxLnBhcmFtKFwidG9rZW5cIikgfSlcbiAgICAgICAgaWYgKCFyZXNldFRva2VuKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yMicsIGBSZXNldCBsaW5rIGlzIGludmFsaWQgb3IgaGFzIGV4cGlyZWRgKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy5yZWRpcmVjdChyZXEub3JpZ2luYWxVcmwpXG4gICAgICAgIH1cblxuICAgICAgICBpZiAocmVzZXRUb2tlbi51c2VkKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yMicsIGBSZXNldCBsaW5rIGlzIGludmFsaWQgb3IgaGFzIGV4cGlyZWRgKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy5yZWRpcmVjdChyZXEub3JpZ2luYWxVcmwpXG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIW1vbWVudChyZXNldFRva2VuLmNyZWF0ZWRfYXQpLmFkZCgyLCBcImRheVwiKS51dGMoKS5pc0FmdGVyKG1vbWVudCgpKSkge1xuICAgICAgICAgICAgcmVxLmZsYXNoKCdlcnJvcjInLCBgUmVzZXQgbGluayBpcyBpbnZhbGlkIG9yIGhhcyBleHBpcmVkYClcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QocmVxLm9yaWdpbmFsVXJsKVxuICAgICAgICB9XG5cbiAgICAgICAgbGV0IGJvZHkgPSByZXEuYm9keVxuXG4gICAgICAgIC8vIHZhbGlkYXRpb25cbiAgICAgICAgaWYgKHZhbGlkYXRvci5pc0VtcHR5KGJvZHkucGFzc3dvcmQpKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgYFBhc3N3b3JkIGlzIHJlcXVpcmVkYCwgJ3Bhc3N3b3JkJywgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QocmVxLm9yaWdpbmFsVXJsKVxuICAgICAgICB9IGVsc2UgaWYgKCF2YWxpZGF0b3IuaXNMZW5ndGgoYm9keS5wYXNzd29yZCwgeyBtaW46IDgsIG1heDogMTIwIH0pKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKHJlcSwgYFBsZWFzZSBlbnRlciBhIHBhc3N3b3JkIGJldHdlZW4gOCBhbmQgMTIwIGNoYXJhY3RlcnMsIHRoZSB0cmlja2llciB0aGUgYmV0dGVyIWAsICdwYXNzd29yZCcsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KHJlcS5vcmlnaW5hbFVybClcbiAgICAgICAgfVxuXG4gICAgICAgIGxldCBhY2NvdW50ID0gYXdhaXQgQWNjb3VudC5nZXQoeyBpZDogcmVzZXRUb2tlbi5hY2NvdW50IH0pXG4gICAgICAgIGlmICghYWNjb3VudCkge1xuICAgICAgICAgICAgcmVxLmZsYXNoKCdlcnJvcjInLCBgUmVzZXQgbGluayBpcyBpbnZhbGlkIG9yIGhhcyBleHBpcmVkJ2ApXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KHJlcS5vcmlnaW5hbFVybClcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIGVuc3VyZSBuZXcgcGFzc3dvcmQgaXMgbm90IHRoZSBzYW1lIGFzIHRoZSBvbGQgb25lXG4gICAgICAgIGlmIChiY3J5cHQuY29tcGFyZVN5bmMoYm9keS5wYXNzd29yZCwgYWNjb3VudC5wYXNzd29yZCkpIHtcbiAgICAgICAgICAgIHNldEZsYXNoRXJyb3IocmVxLCBgTmV3IHBhc3N3b3JkIGNhbm5vdCBiZSBzaW1pbGFyIHRvIGEgcHJldmlvdXNseSB1c2VkIHBhc3N3b3JkYCwgJ3Bhc3N3b3JkJywgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QocmVxLm9yaWdpbmFsVXJsKVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gc2V0IHJlc2V0IHRva2VuIGFzIHVzZWRcbiAgICAgICAgcmVzZXRUb2tlbi5zZXQoJ3VzZWQnLCB0cnVlKVxuICAgICAgICBhd2FpdCByZXNldFRva2VuLnNhdmUoKVxuXG4gICAgICAgIC8vIHVwZGF0ZSBwYXNzd29yZFxuICAgICAgICBsZXQgbmV3UGFzc3dvcmRIYXNoID0gYmNyeXB0Lmhhc2hTeW5jKGJvZHkucGFzc3dvcmQsIDEwKVxuICAgICAgICBhY2NvdW50LnNldCgncGFzc3dvcmQnLCBuZXdQYXNzd29yZEhhc2gpXG4gICAgICAgIGF3YWl0IGFjY291bnQuc2F2ZSgpXG5cbiAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3Jlc2V0X3Bhc3N3b3JkJywgeyBtc2c6ICdQYXNzd29yZCBzdWNjZXNzZnVsbHkgY2hhbmdlZCcgfSlcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn1cblxuLyoqXG4gKiBWZXJpZnkgYW4gYWNjb3VudFxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVxdWVzdH0gcmVxIFRoZSByZXF1ZXN0IG9iamVjdFxuICogQHBhcmFtIHtleHByZXNzLlJlc3BvbnNlfSByZXMgVGhlIHJlc3BvbnNlIG9iamVjdFxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gdmVyaWZ5KHJlcTogYW55LCByZXM6IGFueSkge1xuICAgIHRyeSB7XG5cbiAgICAgICAgbGV0IGNvZGUgPSByZXEucGFyYW0oXCJjb2RlXCIpXG4gICAgICAgIGlmICghY29kZSkge1xuICAgICAgICAgICAgcmV0dXJuIHJlcy5zZXJ2ZXJFcnJvcigpXG4gICAgICAgIH1cblxuICAgICAgICAvLyBnZXQgYWNjb3VudCBieSBjb25maXJtYXRpb24gY29kZVxuICAgICAgICBsZXQgYWNjb3VudCA9IGF3YWl0IEFjY291bnQuZ2V0KHsgY29uZmlybWF0aW9uX2NvZGU6IGNvZGUgfSlcbiAgICAgICAgaWYgKCFhY2NvdW50KSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yJywgJ1NvcnJ5LCBUaGlzIGNvbmZpcm1hdGlvbiBsaW5rIGlzIGludmFsaWQnKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3ZlcmlmeScpXG4gICAgICAgIH1cblxuICAgICAgICBpZiAoYWNjb3VudC5jb25maXJtZWQpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnbXNnJywgJ1RoaXMgYWNjb3VudCBoYXMgYWxyZWFkeSBiZWVuIGNvbmZpcm1lZCcpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvdmVyaWZ5JylcbiAgICAgICAgfVxuXG4gICAgICAgIGFjY291bnQuc2V0KCdjb25maXJtZWQnLCB0cnVlKVxuICAgICAgICBhd2FpdCBhY2NvdW50LnNhdmUoKVxuICAgICAgICByZXEuZmxhc2goJ21zZycsICdHcmVhdCEgWW91ciBhY2NvdW50IGhhcyBiZWVuIGNvbmZpcm1lZCcpXG4gICAgICAgIHJldHVybiByZXMudmlldygnYWNjb3VudC92ZXJpZnknKVxuXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zb2xlLmVycm9yKGUpXG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn0iXX0=
