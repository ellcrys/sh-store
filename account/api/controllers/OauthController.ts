import validator = require('validator')
import Account = require('../models/Account')
import Contract = require('../models/Contract')
import jwt = require('jsonwebtoken')
import moment = require('moment')
import bcrypt = require('bcrypt')

const supportedGrantType = ["client_credentials", "session_token"]

/**
 * Token types
 * @enum {string}
 */
enum TokenType {
    App = "app",
    Session = "session"
}

/**
 * Create an access token
 * 
 * @param {{ id: string, type: string, iat: number, exp?: number }} payload 
 * @returns {string} 
 */
function makeToken(payload: { id: string, type: string, iat: number, exp?: number }): string {
    return jwt.sign(payload, sails.config.app.signingSecret)
}

/**
 * Create an app token
 * 
 * @param {any} res Response object
 * @param {string} clientID The client id
 * @param {string} clientSecret The client secret
 * @returns
 */
async function createAppToken(res: any, clientID: string, clientSecret: string) {
        
    if (validator.isEmpty(clientID)) {
        return res.status(400).json(JSONAPIService.error("400", "Client id is required", "invalid_parameter"))
    }
    if (validator.isEmpty(clientSecret)) {
        return res.status(400).json(JSONAPIService.error("400", "Client secret is required", "invalid_parameter"))
    }
    
    // find contract by client id
    let contract = await Contract.get({ client_id: clientID })
    if (!contract || contract.client_secret != clientSecret) {
        return res.status(401).json(JSONAPIService.error("401", "Client id and client secret are not valid", "invalid_parameter"))
    }
    
    return res.status(201).json({
        access_token: makeToken({ client_id: clientID, type: TokenType.App.toString(), iat: moment().unix() }),
        token_type: TokenType.App.toString()
    })
}

/**
 * Create a session token
 * 
 * @param {*} res 
 * @param {string} email The account email
 * @param {string} password The account password
 * @returns
 */
async function createSessionToken(res: any, email: string, password: string) {
    
    if (validator.isEmpty(email)) {
        return res.status(400).json(JSONAPIService.error("400", "email is required", "invalid_parameter"))
    } else if (!validator.isEmail(email)) {
        return res.status(400).json(JSONAPIService.error("400", "email appears to be invalid", "invalid_parameter"))
    } else if (validator.isEmpty(password)) {
        return res.status(400).json(JSONAPIService.error("400", "password is required", "invalid_parameter"))
    }
    
    // find account by email
    let account = await Account.get({ email })
    if (!account || !bcrypt.compareSync(password, account.password)) {
        return res.status(401).json(JSONAPIService.error("401", "email and password are not valid", "invalid_parameter"))
    }
    
    // unconfirmed account can't be operated on
    if (!account.confirmed) {
        return res.status(403).json(JSONAPIService.error("403", "account has not been confirmed", "unconfirmed_account"))
    }
    
    const exp = moment().add(1, "month").unix()
    return res.status(201).json({
        access_token: makeToken({ id: account.id, type: TokenType.Session.toString(), iat: moment().unix(), exp  }),
        token_type: TokenType.Session.toString(),
        exp,
    })
}

/**
 * Create access token
 * @export
 * @param {Request} req 
 * @param {Response} res 
 */
export async function createToken(req: any, res: any) {
    
    let query = req.query
    let body = req.body
    
    switch (query.grant_type || "") {
        case "":
            return res.status(400).json(JSONAPIService.error("400", "grant type is required", "invalid_parameter"))
        case "client_credentials":
            return createAppToken(res, query.client_id || "", query.client_secret || "")
        case "session":
            return createSessionToken(res, body.email || "", body.password || "")
        default:
            return res.status(400).json(JSONAPIService.error("400", "grant type not supported", "invalid_parameter"))
    }
}

/**
 * User Authorization Page
 * @param req 
 * @param res 
 */
export function authorize(req: any, res: any) {
    return res.send("authorization page")
}