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
const Contract = require("../models/Contract");
const jwt = require("jsonwebtoken");
const moment = require("moment");
const bcrypt = require("bcrypt");
const supportedGrantType = ["client_credentials", "session_token"];
/**
 * Token types
 * @enum {string}
 */
var TokenType;
(function (TokenType) {
    TokenType["App"] = "app";
    TokenType["Session"] = "session";
})(TokenType || (TokenType = {}));
/**
 * Create an access token
 *
 * @param {{ id: string, type: string, iat: number, exp?: number }} payload
 * @returns {string}
 */
function makeToken(payload) {
    return jwt.sign(payload, sails.config.app.signingSecret);
}
/**
 * Create an app token
 *
 * @param {any} res Response object
 * @param {string} clientID The client id
 * @param {string} clientSecret The client secret
 * @returns
 */
function createAppToken(res, clientID, clientSecret) {
    return __awaiter(this, void 0, void 0, function* () {
        if (validator.isEmpty(clientID)) {
            return res.status(400).json(JSONAPIService.error("400", "Client id is required", "invalid_parameter"));
        }
        if (validator.isEmpty(clientSecret)) {
            return res.status(400).json(JSONAPIService.error("400", "Client secret is required", "invalid_parameter"));
        }
        // find contract by client id
        let contract = yield Contract.get({ client_id: clientID });
        if (!contract || contract.client_secret != clientSecret) {
            return res.status(401).json(JSONAPIService.error("401", "Client id and client secret are not valid", "invalid_parameter"));
        }
        return res.status(201).json({
            access_token: makeToken({ client_id: clientID, type: TokenType.App.toString(), iat: moment().unix() }),
            token_type: TokenType.App.toString()
        });
    });
}
/**
 * Create a session token
 *
 * @param {*} res
 * @param {string} email The account email
 * @param {string} password The account password
 * @returns
 */
function createSessionToken(res, email, password) {
    return __awaiter(this, void 0, void 0, function* () {
        if (validator.isEmpty(email)) {
            return res.status(400).json(JSONAPIService.error("400", "email is required", "invalid_parameter"));
        }
        else if (!validator.isEmail(email)) {
            return res.status(400).json(JSONAPIService.error("400", "email appears to be invalid", "invalid_parameter"));
        }
        else if (validator.isEmpty(password)) {
            return res.status(400).json(JSONAPIService.error("400", "password is required", "invalid_parameter"));
        }
        // find account by email
        let account = yield Account.get({ email });
        if (!account || !bcrypt.compareSync(password, account.password)) {
            return res.status(401).json(JSONAPIService.error("401", "email and password are not valid", "invalid_parameter"));
        }
        // unconfirmed account can't be operated on
        if (!account.confirmed) {
            return res.status(403).json(JSONAPIService.error("403", "account has not been confirmed", "unconfirmed_account"));
        }
        const exp = moment().add(1, "month").unix();
        return res.status(201).json({
            access_token: makeToken({ id: account.id, type: TokenType.Session.toString(), iat: moment().unix(), exp }),
            token_type: TokenType.Session.toString(),
            exp,
        });
    });
}
/**
 * Create access token
 * @export
 * @param {Request} req
 * @param {Response} res
 */
function createToken(req, res) {
    return __awaiter(this, void 0, void 0, function* () {
        let query = req.query;
        let body = req.body;
        switch (query.grant_type || "") {
            case "":
                return res.status(400).json(JSONAPIService.error("400", "grant type is required", "invalid_parameter"));
            case "client_credentials":
                return createAppToken(res, query.client_id || "", query.client_secret || "");
            case "session":
                return createSessionToken(res, body.email || "", body.password || "");
            default:
                return res.status(400).json(JSONAPIService.error("400", "grant type not supported", "invalid_parameter"));
        }
    });
}
exports.createToken = createToken;
/**
 * User Authorization Page
 * @param req
 * @param res
 */
function authorize(req, res) {
    return res.send("authorization page");
}
exports.authorize = authorize;

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9jb250cm9sbGVycy9PYXV0aENvbnRyb2xsZXIudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7OztBQUFBLHVDQUF1QztBQUN2Qyw2Q0FBNkM7QUFDN0MsK0NBQStDO0FBQy9DLG9DQUFvQztBQUNwQyxpQ0FBaUM7QUFDakMsaUNBQWlDO0FBRWpDLE1BQU0sa0JBQWtCLEdBQUcsQ0FBQyxvQkFBb0IsRUFBRSxlQUFlLENBQUMsQ0FBQTtBQUVsRTs7O0dBR0c7QUFDSCxJQUFLLFNBR0o7QUFIRCxXQUFLLFNBQVM7SUFDVix3QkFBVyxDQUFBO0lBQ1gsZ0NBQW1CLENBQUE7QUFDdkIsQ0FBQyxFQUhJLFNBQVMsS0FBVCxTQUFTLFFBR2I7QUFFRDs7Ozs7R0FLRztBQUNILG1CQUFtQixPQUFnRTtJQUMvRSxNQUFNLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxPQUFPLEVBQUUsS0FBSyxDQUFDLE1BQU0sQ0FBQyxHQUFHLENBQUMsYUFBYSxDQUFDLENBQUE7QUFDNUQsQ0FBQztBQUVEOzs7Ozs7O0dBT0c7QUFDSCx3QkFBOEIsR0FBUSxFQUFFLFFBQWdCLEVBQUUsWUFBb0I7O1FBRTFFLEVBQUUsQ0FBQyxDQUFDLFNBQVMsQ0FBQyxPQUFPLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDO1lBQzlCLE1BQU0sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLElBQUksQ0FBQyxjQUFjLENBQUMsS0FBSyxDQUFDLEtBQUssRUFBRSx1QkFBdUIsRUFBRSxtQkFBbUIsQ0FBQyxDQUFDLENBQUE7UUFDMUcsQ0FBQztRQUNELEVBQUUsQ0FBQyxDQUFDLFNBQVMsQ0FBQyxPQUFPLENBQUMsWUFBWSxDQUFDLENBQUMsQ0FBQyxDQUFDO1lBQ2xDLE1BQU0sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLElBQUksQ0FBQyxjQUFjLENBQUMsS0FBSyxDQUFDLEtBQUssRUFBRSwyQkFBMkIsRUFBRSxt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a0JBQWtCLENBQUMsR0FBRyxFQUFFLElBQUksQ0FBQyxLQUFLLElBQUksRUFBRSxFQUFFLElBQUksQ0FBQyxRQUFRLElBQUksRUFBRSxDQUFDLENBQUE7WUFDekU7Z0JBQ0ksTUFBTSxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLENBQUMsSUFBSSxDQUFDLGNBQWMsQ0FBQyxLQUFLLENBQUMsS0FBSyxFQUFFLDBCQUEwQixFQUFFLG1CQUFtQixDQUFDLENBQUMsQ0FBQTtRQUNqSCxDQUFDO0lBQ0wsQ0FBQztDQUFBO0FBZkQsa0NBZUM7QUFFRDs7OztHQUlHO0FBQ0gsbUJBQTBCLEdBQVEsRUFBRSxHQUFRO0lBQ3hDLE1BQU0sQ0FBQyxHQUFHLENBQUMsSUFBSSxDQUFDLG9CQUFvQixDQUFDLENBQUE7QUFDekMsQ0FBQztBQUZELDhCQUVDIiwiZmlsZSI6ImFwaS9jb250cm9sbGVycy9PYXV0aENvbnRyb2xsZXIuanMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgdmFsaWRhdG9yID0gcmVxdWlyZSgndmFsaWRhdG9yJylcbmltcG9ydCBBY2NvdW50ID0gcmVxdWlyZSgnLi4vbW9kZWxzL0FjY291bnQnKVxuaW1wb3J0IENvbnRyYWN0ID0gcmVxdWlyZSgnLi4vbW9kZWxzL0NvbnRyYWN0JylcbmltcG9ydCBqd3QgPSByZXF1aXJlKCdqc29ud2VidG9rZW4nKVxuaW1wb3J0IG1vbWVudCA9IHJlcXVpcmUoJ21vbWVudCcpXG5pbXBvcnQgYmNyeXB0ID0gcmVxdWlyZSgnYmNyeXB0JylcblxuY29uc3Qgc3VwcG9ydGVkR3JhbnRUeXBlID0gW1wiY2xpZW50X2NyZWRlbnRpYWxzXCIsIFwic2Vzc2lvbl90b2tlblwiXVxuXG4vKipcbiAqIFRva2VuIHR5cGVzXG4gKiBAZW51bSB7c3RyaW5nfVxuICovXG5lbnVtIFRva2VuVHlwZSB7XG4gICAgQXBwID0gXCJhcHBcIixcbiAgICBTZXNzaW9uID0gXCJzZXNzaW9uXCJcbn1cblxuLyoqXG4gKiBDcmVhdGUgYW4gYWNjZXNzIHRva2VuXG4gKiBcbiAqIEBwYXJhbSB7eyBpZDogc3RyaW5nLCB0eXBlOiBzdHJpbmcsIGlhdDogbnVtYmVyLCBleHA/OiBudW1iZXIgfX0gcGF5bG9hZCBcbiAqIEByZXR1cm5zIHtzdHJpbmd9IFxuICovXG5mdW5jdGlvbiBtYWtlVG9rZW4ocGF5bG9hZDogeyBpZDogc3RyaW5nLCB0eXBlOiBzdHJpbmcsIGlhdDogbnVtYmVyLCBleHA/OiBudW1iZXIgfSk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGp3dC5zaWduKHBheWxvYWQsIHNhaWxzLmNvbmZpZy5hcHAuc2lnbmluZ1NlY3JldClcbn1cblxuLyoqXG4gKiBDcmVhdGUgYW4gYXBwIHRva2VuXG4gKiBcbiAqIEBwYXJhbSB7YW55fSByZXMgUmVzcG9uc2Ugb2JqZWN0XG4gKiBAcGFyYW0ge3N0cmluZ30gY2xpZW50SUQgVGhlIGNsaWVudCBpZFxuICogQHBhcmFtIHtzdHJpbmd9IGNsaWVudFNlY3JldCBUaGUgY2xpZW50IHNlY3JldFxuICogQHJldHVybnNcbiAqL1xuYXN5bmMgZnVuY3Rpb24gY3JlYXRlQXBwVG9rZW4ocmVzOiBhbnksIGNsaWVudElEOiBzdHJpbmcsIGNsaWVudFNlY3JldDogc3RyaW5nKSB7XG4gICAgICAgIFxuICAgIGlmICh2YWxpZGF0b3IuaXNFbXB0eShjbGllbnRJRCkpIHtcbiAgICAgICAgcmV0dXJuIHJlcy5zdGF0dXMoNDAwKS5qc29uKEpTT05BUElTZXJ2aWNlLmVycm9yKFwiNDAwXCIsIFwiQ2xpZW50IGlkIGlzIHJlcXVpcmVkXCIsIFwiaW52YWxpZF9wYXJhbWV0ZXJcIikpXG4gICAgfVxuICAgIGlmICh2YWxpZGF0b3IuaXNFbXB0eShjbGllbnRTZWNyZXQpKSB7XG4gICAgICAgIHJldHVybiByZXMuc3RhdHVzKDQwMCkuanNvbihKU09OQVBJU2VydmljZS5lcnJvcihcIjQwMFwiLCBcIkNsaWVudCBzZWNyZXQgaXMgcmVxdWlyZWRcIiwgXCJpbnZhbGlkX3BhcmFtZXRlclwiKSlcbiAgICB9XG4gICAgXG4gICAgLy8gZmluZCBjb250cmFjdCBieSBjbGllbnQgaWRcbiAgICBsZXQgY29udHJhY3QgPSBhd2FpdCBDb250cmFjdC5nZXQoeyBjbGllbnRfaWQ6IGNsaWVudElEIH0pXG4gICAgaWYgKCFjb250cmFjdCB8fCBjb250cmFjdC5jbGllbnRfc2VjcmV0ICE9IGNsaWVudFNlY3JldCkge1xuICAgICAgICByZXR1cm4gcmVzLnN0YXR1cyg0MDEpLmpzb24oSlNPTkFQSVNlcnZpY2UuZXJyb3IoXCI0MDFcIiwgXCJDbGllbnQgaWQgYW5kIGNsaWVudCBzZWNyZXQgYXJlIG5vdCB2YWxpZFwiLCBcImludmFsaWRfcGFyYW1ldGVyXCIpKVxuICAgIH1cbiAgICBcbiAgICByZXR1cm4gcmVzLnN0YXR1cygyMDEpLmpzb24oe1xuICAgICAgICBhY2Nlc3NfdG9rZW46IG1ha2VUb2tlbih7IGNsaWVudF9pZDogY2xpZW50SUQsIHR5cGU6IFRva2VuVHlwZS5BcHAudG9TdHJpbmcoKSwgaWF0OiBtb21lbnQoKS51bml4KCkgfSksXG4gICAgICAgIHRva2VuX3R5cGU6IFRva2VuVHlwZS5BcHAudG9TdHJpbmcoKVxuICAgIH0pXG59XG5cbi8qKlxuICogQ3JlYXRlIGEgc2Vzc2lvbiB0b2tlblxuICogXG4gKiBAcGFyYW0geyp9IHJlcyBcbiAqIEBwYXJhbSB7c3RyaW5nfSBlbWFpbCBUaGUgYWNjb3VudCBlbWFpbFxuICogQHBhcmFtIHtzdHJpbmd9IHBhc3N3b3JkIFRoZSBhY2NvdW50IHBhc3N3b3JkXG4gKiBAcmV0dXJuc1xuICovXG5hc3luYyBmdW5jdGlvbiBjcmVhdGVTZXNzaW9uVG9rZW4ocmVzOiBhbnksIGVtYWlsOiBzdHJpbmcsIHBhc3N3b3JkOiBzdHJpbmcpIHtcbiAgICBcbiAgICBpZiAodmFsaWRhdG9yLmlzRW1wdHkoZW1haWwpKSB7XG4gICAgICAgIHJldHVybiByZXMuc3RhdHVzKDQwMCkuanNvbihKU09OQVBJU2VydmljZS5lcnJvcihcIjQwMFwiLCBcImVtYWlsIGlzIHJlcXVpcmVkXCIsIFwiaW52YWxpZF9wYXJhbWV0ZXJcIikpXG4gICAgfSBlbHNlIGlmICghdmFsaWRhdG9yLmlzRW1haWwoZW1haWwpKSB7XG4gICAgICAgIHJldHVybiByZXMuc3RhdHVzKDQwMCkuanNvbihKU09OQVBJU2VydmljZS5lcnJvcihcIjQwMFwiLCBcImVtYWlsIGFwcGVhcnMgdG8gYmUgaW52YWxpZFwiLCBcImludmFsaWRfcGFyYW1ldGVyXCIpKVxuICAgIH0gZWxzZSBpZiAodmFsaWRhdG9yLmlzRW1wdHkocGFzc3dvcmQpKSB7XG4gICAgICAgIHJldHVybiByZXMuc3RhdHVzKDQwMCkuanNvbihKU09OQVBJU2VydmljZS5lcnJvcihcIjQwMFwiLCBcInBhc3N3b3JkIGlzIHJlcXVpcmVkXCIsIFwiaW52YWxpZF9wYXJhbWV0ZXJcIikpXG4gICAgfVxuICAgIFxuICAgIC8vIGZpbmQgYWNjb3VudCBieSBlbWFpbFxuICAgIGxldCBhY2NvdW50ID0gYXdhaXQgQWNjb3VudC5nZXQoeyBlbWFpbCB9KVxuICAgIGlmICghYWNjb3VudCB8fCAhYmNyeXB0LmNvbXBhcmVTeW5jKHBhc3N3b3JkLCBhY2NvdW50LnBhc3N3b3JkKSkge1xuICAgICAgICByZXR1cm4gcmVzLnN0YXR1cyg0MDEpLmpzb24oSlNPTkFQSVNlcnZpY2UuZXJyb3IoXCI0MDFcIiwgXCJlbWFpbCBhbmQgcGFzc3dvcmQgYXJlIG5vdCB2YWxpZFwiLCBcImludmFsaWRfcGFyYW1ldGVyXCIpKVxuICAgIH1cbiAgICBcbiAgICAvLyB1bmNvbmZpcm1lZCBhY2NvdW50IGNhbid0IGJlIG9wZXJhdGVkIG9uXG4gICAgaWYgKCFhY2NvdW50LmNvbmZpcm1lZCkge1xuICAgICAgICByZXR1cm4gcmVzLnN0YXR1cyg0MDMpLmpzb24oSlNPTkFQSVNlcnZpY2UuZXJyb3IoXCI0MDNcIiwgXCJhY2NvdW50IGhhcyBub3QgYmVlbiBjb25maXJtZWRcIiwgXCJ1bmNvbmZpcm1lZF9hY2NvdW50XCIpKVxuICAgIH1cbiAgICBcbiAgICBjb25zdCBleHAgPSBtb21lbnQoKS5hZGQoMSwgXCJtb250aFwiKS51bml4KClcbiAgICByZXR1cm4gcmVzLnN0YXR1cygyMDEpLmpzb24oe1xuICAgICAgICBhY2Nlc3NfdG9rZW46IG1ha2VUb2tlbih7IGlkOiBhY2NvdW50LmlkLCB0eXBlOiBUb2tlblR5cGUuU2Vzc2lvbi50b1N0cmluZygpLCBpYXQ6IG1vbWVudCgpLnVuaXgoKSwgZXhwICB9KSxcbiAgICAgICAgdG9rZW5fdHlwZTogVG9rZW5UeXBlLlNlc3Npb24udG9TdHJpbmcoKSxcbiAgICAgICAgZXhwLFxuICAgIH0pXG59XG5cbi8qKlxuICogQ3JlYXRlIGFjY2VzcyB0b2tlblxuICogQGV4cG9ydFxuICogQHBhcmFtIHtSZXF1ZXN0fSByZXEgXG4gKiBAcGFyYW0ge1Jlc3BvbnNlfSByZXMgXG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBjcmVhdGVUb2tlbihyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICBcbiAgICBsZXQgcXVlcnkgPSByZXEucXVlcnlcbiAgICBsZXQgYm9keSA9IHJlcS5ib2R5XG4gICAgXG4gICAgc3dpdGNoIChxdWVyeS5ncmFudF90eXBlIHx8IFwiXCIpIHtcbiAgICAgICAgY2FzZSBcIlwiOlxuICAgICAgICAgICAgcmV0dXJuIHJlcy5zdGF0dXMoNDAwKS5qc29uKEpTT05BUElTZXJ2aWNlLmVycm9yKFwiNDAwXCIsIFwiZ3JhbnQgdHlwZSBpcyByZXF1aXJlZFwiLCBcImludmFsaWRfcGFyYW1ldGVyXCIpKVxuICAgICAgICBjYXNlIFwiY2xpZW50X2NyZWRlbnRpYWxzXCI6XG4gICAgICAgICAgICByZXR1cm4gY3JlYXRlQXBwVG9rZW4ocmVzLCBxdWVyeS5jbGllbnRfaWQgfHwgXCJcIiwgcXVlcnkuY2xpZW50X3NlY3JldCB8fCBcIlwiKVxuICAgICAgICBjYXNlIFwic2Vzc2lvblwiOlxuICAgICAgICAgICAgcmV0dXJuIGNyZWF0ZVNlc3Npb25Ub2tlbihyZXMsIGJvZHkuZW1haWwgfHwgXCJcIiwgYm9keS5wYXNzd29yZCB8fCBcIlwiKVxuICAgICAgICBkZWZhdWx0OlxuICAgICAgICAgICAgcmV0dXJuIHJlcy5zdGF0dXMoNDAwKS5qc29uKEpTT05BUElTZXJ2aWNlLmVycm9yKFwiNDAwXCIsIFwiZ3JhbnQgdHlwZSBub3Qgc3VwcG9ydGVkXCIsIFwiaW52YWxpZF9wYXJhbWV0ZXJcIikpXG4gICAgfVxufVxuXG4vKipcbiAqIFVzZXIgQXV0aG9yaXphdGlvbiBQYWdlXG4gKiBAcGFyYW0gcmVxIFxuICogQHBhcmFtIHJlcyBcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGF1dGhvcml6ZShyZXE6IGFueSwgcmVzOiBhbnkpIHtcbiAgICByZXR1cm4gcmVzLnNlbmQoXCJhdXRob3JpemF0aW9uIHBhZ2VcIilcbn0iXX0=
