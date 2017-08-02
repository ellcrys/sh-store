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
const uuid4 = require("uuid4");
const random = require("randomstring");
const Random = require("random-js");
const bcrypt = require("bcrypt");
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
            function setFlashError(msg, field, body) {
                req.flash("error", msg);
                req.flash("field", field);
                req.flash("form", body);
            }
            // send validation error
            if (valError != null) {
                setFlashError(valError.msg, valError.field, body);
                return res.redirect('back');
            }
            // check if an account with matching email exists
            let existingAccount = yield Account.get({ email: body.email });
            if (existingAccount) {
                setFlashError("Email has already been registered", "email", body);
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
            function setFlashError(msg, field, body) {
                req.flash("error", msg);
                req.flash("field", field);
                req.flash("form", body);
            }
            if (!body.email) {
                setFlashError('Please enter your email', 'email', body);
                return res.redirect('back');
            }
            if (!body.password) {
                setFlashError('Please enter your password', 'password', body);
                return res.redirect('back');
            }
            let account = yield Account.get({ email: body.email });
            if (!account) {
                setFlashError('Your email and password combination is not valid', '', body);
                return res.redirect('back');
            }
            if (!bcrypt.compareSync(body.password, account.password)) {
                setFlashError('Your email and password combination is not valid', '', body);
                return res.redirect('back');
            }
            if (!account.confirmed) {
                req.flash('warn', 'warn');
                setFlashError(`Your account has not been confirmed. Please check your email for the confirmation message. <a href="/account/resend_code/${account.id}">Resend</a>`, '', body);
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

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9jb250cm9sbGVycy9BY2NvdW50Q29udHJvbGxlci50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7O0FBRUEsdUNBQXdDO0FBQ3hDLDZDQUE4QztBQUM5QywrQkFBZ0M7QUFDaEMsdUNBQXdDO0FBQ3hDLG9DQUFxQztBQUNyQyxpQ0FBa0M7QUFFbEM7Ozs7O0dBS0c7QUFDSCx5QkFBZ0MsT0FBZ0I7SUFDNUMsRUFBRSxDQUFDLENBQUMsU0FBUyxDQUFDLE9BQU8sQ0FBQyxPQUFPLENBQUMsVUFBVSxDQUFDLENBQUM7UUFBQyxNQUFNLENBQUMsRUFBRSxHQUFHLEVBQUUsd0JBQXdCLEVBQUUsS0FBSyxFQUFFLFlBQVksRUFBRSxDQUFDO0lBQ3pHLEVBQUUsQ0FBQyxDQUFDLFNBQVMsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLFNBQVMsQ0FBQyxDQUFDO1FBQUMsTUFBTSxDQUFDLEVBQUUsR0FBRyxFQUFFLHVCQUF1QixFQUFFLEtBQUssRUFBRSxXQUFXLEVBQUUsQ0FBQztJQUN0RyxFQUFFLENBQUMsQ0FBQyxTQUFTLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxLQUFLLENBQUMsQ0FBQztRQUFDLE1BQU0sQ0FBQyxFQUFFLEdBQUcsRUFBRSxtQkFBbUIsRUFBRSxLQUFLLEVBQUUsT0FBTyxFQUFFLENBQUM7SUFDMUYsRUFBRSxDQUFDLENBQUMsQ0FBQyxTQUFTLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxLQUFLLENBQUMsQ0FBQztRQUFDLE1BQU0sQ0FBQyxFQUFFLEdBQUcsRUFBRSxtQ0FBbUMsRUFBRSxLQUFLLEVBQUUsT0FBTyxFQUFFLENBQUM7SUFDM0csRUFBRSxDQUFDLENBQUMsU0FBUyxDQUFDLE9BQU8sQ0FBQyxPQUFPLENBQUMsUUFBUSxDQUFDLENBQUM7UUFBQyxNQUFNLENBQUMsRUFBRSxHQUFHLEVBQUUsc0JBQXNCLEVBQUUsS0FBSyxFQUFFLFVBQVUsRUFBRSxDQUFDO0lBQ25HLEVBQUUsQ0FBQyxDQUFDLENB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ImFwaS9jb250cm9sbGVycy9BY2NvdW50Q29udHJvbGxlci5qcyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IGFjY2VzcyB9IGZyb20gJ2ZzJztcbmltcG9ydCBleHByZXNzID0gcmVxdWlyZSgnZXhwcmVzcycpO1xuaW1wb3J0IHZhbGlkYXRvciA9IHJlcXVpcmUoJ3ZhbGlkYXRvcicpO1xuaW1wb3J0IEFjY291bnQgPSByZXF1aXJlKCcuLi9tb2RlbHMvQWNjb3VudCcpO1xuaW1wb3J0IHV1aWQ0ID0gcmVxdWlyZSgndXVpZDQnKTtcbmltcG9ydCByYW5kb20gPSByZXF1aXJlKCdyYW5kb21zdHJpbmcnKTtcbmltcG9ydCBSYW5kb20gPSByZXF1aXJlKCdyYW5kb20tanMnKTtcbmltcG9ydCBiY3J5cHQgPSByZXF1aXJlKCdiY3J5cHQnKTtcblxuLyoqXG4gKiBWYWxpZGF0ZXMgcmVzLmJvZHkgY29udGFpbmluZyBhY2NvdW50IGRldGFpbHNcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7QWNjb3VudH0gYWNjb3VudCBcbiAqIEByZXR1cm5zIHtWYWxpZGF0aW9uRXJyb3J9IFxuICovXG5leHBvcnQgZnVuY3Rpb24gdmFsaWRhdGVBY2NvdW50KGFjY291bnQ6IEFjY291bnQpOiBWYWxpZGF0aW9uRXJyb3Ige1xuICAgIGlmICh2YWxpZGF0b3IuaXNFbXB0eShhY2NvdW50LmZpcnN0X25hbWUpKSByZXR1cm4geyBtc2c6IFwiRmlyc3QgbmFtZSBpcyByZXF1aXJlZFwiLCBmaWVsZDogXCJmaXJzdF9uYW1lXCIgfTtcbiAgICBpZiAodmFsaWRhdG9yLmlzRW1wdHkoYWNjb3VudC5sYXN0X25hbWUpKSByZXR1cm4geyBtc2c6IFwiTGFzdCBuYW1lIGlzIHJlcXVpcmVkXCIsIGZpZWxkOiBcImxhc3RfbmFtZVwiIH07XG4gICAgaWYgKHZhbGlkYXRvci5pc0VtcHR5KGFjY291bnQuZW1haWwpKSByZXR1cm4geyBtc2c6IFwiRW1haWwgaXMgcmVxdWlyZWRcIiwgZmllbGQ6IFwiZW1haWxcIiB9O1xuICAgIGlmICghdmFsaWRhdG9yLmlzRW1haWwoYWNjb3VudC5lbWFpbCkpIHJldHVybiB7IG1zZzogXCJFbWFpbCBkb2VzIG5vdCBhcHBlYXIgdG8gYmUgdmFsaWRcIiwgZmllbGQ6IFwiZW1haWxcIiB9O1xuICAgIGlmICh2YWxpZGF0b3IuaXNFbXB0eShhY2NvdW50LnBhc3N3b3JkKSkgcmV0dXJuIHsgbXNnOiBcIlBhc3N3b3JkIGlzIHJlcXVpcmVkXCIsIGZpZWxkOiBcInBhc3N3b3JkXCIgfTtcbiAgICBpZiAoIXZhbGlkYXRvci5pc0xlbmd0aChhY2NvdW50LnBhc3N3b3JkLCB7IG1pbjogOCwgbWF4OiAxMjAgfSkpIHJldHVybiB7IG1zZzogXCJQbGVhc2UgZW50ZXIgYSBwYXNzd29yZCBiZXR3ZWVuIDggYW5kIDEyMCBjaGFyYWN0ZXJzLCB0aGUgdHJpY2tpZXIgdGhlIGJldHRlciFcIiwgZmllbGQ6IFwicGFzc3dvcmRcIiB9O1xuICAgIHJldHVybiBudWxsXG59XG5cbi8qKlxuICogSGFuZGxlIHNpZ251cCBwcm9jZXNzaW5nXG4gKiBcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7ZXhwcmVzcy5SZXF1ZXN0fSByZXEgVGhlIHJlcXVlc3Qgb2JqZWN0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVzcG9uc2V9IHJlcyBUaGUgcmVzcG9uc2Ugb2JqZWN0XG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBzaWdudXAocmVxOiBhbnksIHJlczogYW55KSB7XG4gICAgdHJ5IHtcblxuICAgICAgICBsZXQgYm9keSA9IChyZXEuYm9keSBhcyBBY2NvdW50KVxuICAgICAgICBjb25zdCB2YWxFcnJvciA9IHZhbGlkYXRlQWNjb3VudChib2R5KVxuXG4gICAgICAgIGZ1bmN0aW9uIHNldEZsYXNoRXJyb3IobXNnOiBzdHJpbmcsIGZpZWxkOiBzdHJpbmcsIGJvZHk6IGFueSkge1xuICAgICAgICAgICAgcmVxLmZsYXNoKFwiZXJyb3JcIiwgbXNnKVxuICAgICAgICAgICAgcmVxLmZsYXNoKFwiZmllbGRcIiwgZmllbGQpXG4gICAgICAgICAgICByZXEuZmxhc2goXCJmb3JtXCIsIGJvZHkpXG4gICAgICAgIH1cblxuICAgICAgICAvLyBzZW5kIHZhbGlkYXRpb24gZXJyb3JcbiAgICAgICAgaWYgKHZhbEVycm9yICE9IG51bGwpIHtcbiAgICAgICAgICAgIHNldEZsYXNoRXJyb3IodmFsRXJyb3IubXNnLCB2YWxFcnJvci5maWVsZCwgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gY2hlY2sgaWYgYW4gYWNjb3VudCB3aXRoIG1hdGNoaW5nIGVtYWlsIGV4aXN0c1xuICAgICAgICBsZXQgZXhpc3RpbmdBY2NvdW50ID0gYXdhaXQgQWNjb3VudC5nZXQoeyBlbWFpbDogYm9keS5lbWFpbCB9KVxuICAgICAgICBpZiAoZXhpc3RpbmdBY2NvdW50KSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKFwiRW1haWwgaGFzIGFscmVhZHkgYmVlbiByZWdpc3RlcmVkXCIsIFwiZW1haWxcIiwgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gc2V0IGlkLCBjbGllbnQgaWQgYW5kIHNlY3JldFxuICAgICAgICBib2R5LmlkID0gdXVpZDQoKVxuICAgICAgICBib2R5LmNsaWVudF9pZCA9IHJhbmRvbS5nZW5lcmF0ZSgzMilcbiAgICAgICAgYm9keS5jbGllbnRfc2VjcmV0ID0gcmFuZG9tLmdlbmVyYXRlKFJhbmRvbS5pbnRlZ2VyKDI3LCAzMikoUmFuZG9tLmVuZ2luZXMubXQxOTkzNygpLmF1dG9TZWVkKCkpKVxuICAgICAgICBib2R5LnBhc3N3b3JkID0gYmNyeXB0Lmhhc2hTeW5jKGJvZHkucGFzc3dvcmQsIDEwKVxuICAgICAgICBib2R5LmNvbmZpcm1hdGlvbl9jb2RlID0gcmFuZG9tLmdlbmVyYXRlKDI4KVxuICAgICAgICBsZXQgYWNjb3VudCA9IGF3YWl0IEFjY291bnQuY3JlYXRlKGJvZHkpXG5cbiAgICAgICAgLy8gc2VuZCBjb25maXJtYXRpb24gbWVzc2FnZVxuICAgICAgICBsZXQgcmVzcCA9IGF3YWl0IEVtYWlsU2VydmljZS5zZW5kQWNjb3VudENvbmZpcm1hdGlvbihib2R5LmVtYWlsLCBhY2NvdW50LnRvSlNPTigpKVxuXG4gICAgICAgIHJlcS5mbGFzaChcImFjY291bnRcIiwgYWNjb3VudClcbiAgICAgICAgcmVzLnJlZGlyZWN0KFwiL2NvbmZpcm1cIilcblxuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihlKVxuICAgICAgICByZXR1cm4gcmVzLnNlcnZlckVycm9yKClcbiAgICB9XG59XG5cbi8qKlxuICogTG9naW4gcHJvY2Vzc2luZ1xuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVxdWVzdH0gcmVxIFRoZSByZXF1ZXN0IG9iamVjdFxuICogQHBhcmFtIHtleHByZXNzLlJlc3BvbnNlfSByZXMgVGhlIHJlc3BvbnNlIG9iamVjdFxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gbG9naW4ocmVxOiBhbnksIHJlczogYW55KSB7XG4gICAgdHJ5IHtcbiAgICAgICAgbGV0IGJvZHkgPSByZXEuYm9keVxuICAgICAgICBcbiAgICAgICAgZnVuY3Rpb24gc2V0Rmxhc2hFcnJvcihtc2c6IHN0cmluZywgZmllbGQ6IHN0cmluZywgYm9keTogYW55KSB7XG4gICAgICAgICAgICByZXEuZmxhc2goXCJlcnJvclwiLCBtc2cpXG4gICAgICAgICAgICByZXEuZmxhc2goXCJmaWVsZFwiLCBmaWVsZClcbiAgICAgICAgICAgIHJlcS5mbGFzaChcImZvcm1cIiwgYm9keSlcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgaWYgKCFib2R5LmVtYWlsKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKCdQbGVhc2UgZW50ZXIgeW91ciBlbWFpbCcsICdlbWFpbCcsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KCdiYWNrJylcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgaWYgKCFib2R5LnBhc3N3b3JkKSB7XG4gICAgICAgICAgICBzZXRGbGFzaEVycm9yKCdQbGVhc2UgZW50ZXIgeW91ciBwYXNzd29yZCcsICdwYXNzd29yZCcsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KCdiYWNrJylcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgbGV0IGFjY291bnQgPSBhd2FpdCBBY2NvdW50LmdldCh7IGVtYWlsOiBib2R5LmVtYWlsIH0pXG4gICAgICAgIGlmICghYWNjb3VudCkge1xuICAgICAgICAgICAgc2V0Rmxhc2hFcnJvcignWW91ciBlbWFpbCBhbmQgcGFzc3dvcmQgY29tYmluYXRpb24gaXMgbm90IHZhbGlkJywgJycsIGJvZHkpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnJlZGlyZWN0KCdiYWNrJylcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgaWYgKCFiY3J5cHQuY29tcGFyZVN5bmMoYm9keS5wYXNzd29yZCwgYWNjb3VudC5wYXNzd29yZCkpIHtcbiAgICAgICAgICAgIHNldEZsYXNoRXJyb3IoJ1lvdXIgZW1haWwgYW5kIHBhc3N3b3JkIGNvbWJpbmF0aW9uIGlzIG5vdCB2YWxpZCcsICcnLCBib2R5KVxuICAgICAgICAgICAgcmV0dXJuIHJlcy5yZWRpcmVjdCgnYmFjaycpXG4gICAgICAgIH1cbiAgICAgICAgXG4gICAgICAgIGlmICghYWNjb3VudC5jb25maXJtZWQpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnd2FybicsICd3YXJuJylcbiAgICAgICAgICAgIHNldEZsYXNoRXJyb3IoYFlvdXIgYWNjb3VudCBoYXMgbm90IGJlZW4gY29uZmlybWVkLiBQbGVhc2UgY2hlY2sgeW91ciBlbWFpbCBmb3IgdGhlIGNvbmZpcm1hdGlvbiBtZXNzYWdlLiA8YSBocmVmPVwiL2FjY291bnQvcmVzZW5kX2NvZGUvJHthY2NvdW50LmlkfVwiPlJlc2VuZDwvYT5gLCAnJywgYm9keSlcbiAgICAgICAgICAgIHJldHVybiByZXMucmVkaXJlY3QoJ2JhY2snKVxuICAgICAgICB9XG4gICAgICAgIFxuICAgICAgICByZXR1cm4gcmVzLnNlbmQoYWNjb3VudClcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn1cblxuLyoqXG4gKiBWZXJpZnkgYW4gYWNjb3VudFxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVxdWVzdH0gcmVxIFRoZSByZXF1ZXN0IG9iamVjdFxuICogQHBhcmFtIHtleHByZXNzLlJlc3BvbnNlfSByZXMgVGhlIHJlc3BvbnNlIG9iamVjdFxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gcmVzZW5kQ29uZmlybWF0aW9uKHJlcTogYW55LCByZXM6IGFueSkge1xuICAgIHRyeSB7XG4gICAgICAgIFxuICAgICAgICBsZXQgdXNlcklEID0gcmVxLnBhcmFtKFwiaWRcIilcbiAgICAgICAgaWYgKCF1c2VySUQpIHtcbiAgICAgICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gZ2V0IGFjY291bnQgYnkgdXNlciBpZFxuICAgICAgICBsZXQgYWNjb3VudCA9IGF3YWl0IEFjY291bnQuZ2V0KHsgaWQ6IHVzZXJJRCB9KVxuICAgICAgICBpZiAoIWFjY291bnQpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnZXJyb3InLCAnU29ycnksIFRoaXMgY29uZmlybWF0aW9uIGxpbmsgaXMgaW52YWxpZCcpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvdmVyaWZ5JylcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgaWYgKGFjY291bnQuY29uZmlybWVkKSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ21zZycsICdUaGlzIGFjY291bnQgaGFzIGFscmVhZHkgYmVlbiBjb25maXJtZWQnKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3ZlcmlmeScpXG4gICAgICAgIH1cbiAgICAgICAgXG4gICAgICAgIC8vIHNlbmQgY29uZmlybWF0aW9uIG1lc3NhZ2VcbiAgICAgICAgbGV0IHJlc3AgPSBhd2FpdCBFbWFpbFNlcnZpY2Uuc2VuZEFjY291bnRDb25maXJtYXRpb24oYWNjb3VudC5lbWFpbCwgYWNjb3VudC50b0pTT04oKSlcbiAgICAgICAgcmVxLmZsYXNoKCdtc2cnLCAnQ29uZmlybWF0aW9uIG1lc3NhZ2UgaGFzIGJlZW4gc2VudC4gUGxlYXNlIGNoZWNrIHlvdXIgZW1haWwnKVxuICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvdmVyaWZ5JylcbiAgICAgICAgXG4gICAgfSBjYXRjaChlKSB7XG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn1cblxuLyoqXG4gKiBWZXJpZnkgYW4gYWNjb3VudFxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge2V4cHJlc3MuUmVxdWVzdH0gcmVxIFRoZSByZXF1ZXN0IG9iamVjdFxuICogQHBhcmFtIHtleHByZXNzLlJlc3BvbnNlfSByZXMgVGhlIHJlc3BvbnNlIG9iamVjdFxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gdmVyaWZ5KHJlcTogYW55LCByZXM6IGFueSkge1xuICAgIHRyeSB7XG5cbiAgICAgICAgbGV0IGNvZGUgPSByZXEucGFyYW0oXCJjb2RlXCIpXG4gICAgICAgIGlmICghY29kZSkge1xuICAgICAgICAgICAgcmV0dXJuIHJlcy5zZXJ2ZXJFcnJvcigpXG4gICAgICAgIH1cblxuICAgICAgICAvLyBnZXQgYWNjb3VudCBieSBjb25maXJtYXRpb24gY29kZVxuICAgICAgICBsZXQgYWNjb3VudCA9IGF3YWl0IEFjY291bnQuZ2V0KHsgY29uZmlybWF0aW9uX2NvZGU6IGNvZGUgfSlcbiAgICAgICAgaWYgKCFhY2NvdW50KSB7XG4gICAgICAgICAgICByZXEuZmxhc2goJ2Vycm9yJywgJ1NvcnJ5LCBUaGlzIGNvbmZpcm1hdGlvbiBsaW5rIGlzIGludmFsaWQnKVxuICAgICAgICAgICAgcmV0dXJuIHJlcy52aWV3KCdhY2NvdW50L3ZlcmlmeScpXG4gICAgICAgIH1cblxuICAgICAgICBpZiAoYWNjb3VudC5jb25maXJtZWQpIHtcbiAgICAgICAgICAgIHJlcS5mbGFzaCgnbXNnJywgJ1RoaXMgYWNjb3VudCBoYXMgYWxyZWFkeSBiZWVuIGNvbmZpcm1lZCcpXG4gICAgICAgICAgICByZXR1cm4gcmVzLnZpZXcoJ2FjY291bnQvdmVyaWZ5JylcbiAgICAgICAgfVxuXG4gICAgICAgIGFjY291bnQuc2V0KCdjb25maXJtZWQnLCB0cnVlKVxuICAgICAgICBhd2FpdCBhY2NvdW50LnNhdmUoKVxuICAgICAgICByZXEuZmxhc2goJ21zZycsICdHcmVhdCEgWW91ciBhY2NvdW50IGhhcyBiZWVuIGNvbmZpcm1lZCcpXG4gICAgICAgIHJldHVybiByZXMudmlldygnYWNjb3VudC92ZXJpZnknKVxuXG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zb2xlLmVycm9yKGUpXG4gICAgICAgIHJldHVybiByZXMuc2VydmVyRXJyb3IoKVxuICAgIH1cbn0iXX0=
