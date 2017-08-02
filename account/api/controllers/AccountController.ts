import { access } from 'fs';
import express = require('express');
import validator = require('validator');
import Account = require('../models/Account');
import uuid4 = require('uuid4');
import random = require('randomstring');
import Random = require('random-js');
import bcrypt = require('bcrypt');

/**
 * Validates res.body containing account details
 * @export
 * @param {Account} account 
 * @returns {ValidationError} 
 */
export function validateAccount(account: Account): ValidationError {
    if (validator.isEmpty(account.first_name)) return { msg: "First name is required", field: "first_name" };
    if (validator.isEmpty(account.last_name)) return { msg: "Last name is required", field: "last_name" };
    if (validator.isEmpty(account.email)) return { msg: "Email is required", field: "email" };
    if (!validator.isEmail(account.email)) return { msg: "Email does not appear to be valid", field: "email" };
    if (validator.isEmpty(account.password)) return { msg: "Password is required", field: "password" };
    if (!validator.isLength(account.password, { min: 8, max: 120 })) return { msg: "Please enter a password between 8 and 120 characters, the trickier the better!", field: "password" };
    return null
}

/**
 * Handle signup processing
 * 
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
export async function signup(req: any, res: any) {
    try {

        let body = (req.body as Account)
        const valError = validateAccount(body)

        function setFlashError(msg: string, field: string, body: any) {
            req.flash("error", msg)
            req.flash("field", field)
            req.flash("form", body)
        }

        // send validation error
        if (valError != null) {
            setFlashError(valError.msg, valError.field, body)
            return res.redirect('back')
        }

        // check if an account with matching email exists
        let existingAccount = await Account.get({ email: body.email })
        if (existingAccount) {
            setFlashError("Email has already been registered", "email", body)
            return res.redirect('back')
        }

        // set id, client id and secret
        body.id = uuid4()
        body.client_id = random.generate(32)
        body.client_secret = random.generate(Random.integer(27, 32)(Random.engines.mt19937().autoSeed()))
        body.password = bcrypt.hashSync(body.password, 10)
        body.confirmation_code = random.generate(28)
        let account = await Account.create(body)

        // send confirmation message
        let resp = await EmailService.sendAccountConfirmation(body.email, account.toJSON())

        req.flash("account", account)
        res.redirect("/confirm")

    } catch (e) {
        console.error(e)
        return res.serverError()
    }
}

/**
 * Login processing
 * 
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
export async function login(req: any, res: any) {
    try {
        let body = req.body
        
        function setFlashError(msg: string, field: string, body: any) {
            req.flash("error", msg)
            req.flash("field", field)
            req.flash("form", body)
        }
        
        if (!body.email) {
            setFlashError('Please enter your email', 'email', body)
            return res.redirect('back')
        }
        
        if (!body.password) {
            setFlashError('Please enter your password', 'password', body)
            return res.redirect('back')
        }
        
        let account = await Account.get({ email: body.email })
        if (!account) {
            setFlashError('Your email and password combination is not valid', '', body)
            return res.redirect('back')
        }
        
        if (!bcrypt.compareSync(body.password, account.password)) {
            setFlashError('Your email and password combination is not valid', '', body)
            return res.redirect('back')
        }
        
        if (!account.confirmed) {
            req.flash('warn', 'warn')
            setFlashError(`Your account has not been confirmed. Please check your email for the confirmation message. <a href="/account/resend_code/${account.id}">Resend</a>`, '', body)
            return res.redirect('back')
        }
        
        return res.send(account)
    } catch (e) {
        return res.serverError()
    }
}

/**
 * Verify an account
 * 
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
export async function resendConfirmation(req: any, res: any) {
    try {
        
        let userID = req.param("id")
        if (!userID) {
            return res.serverError()
        }

        // get account by user id
        let account = await Account.get({ id: userID })
        if (!account) {
            req.flash('error', 'Sorry, This confirmation link is invalid')
            return res.view('account/verify')
        }
        
        if (account.confirmed) {
            req.flash('msg', 'This account has already been confirmed')
            return res.view('account/verify')
        }
        
        // send confirmation message
        let resp = await EmailService.sendAccountConfirmation(account.email, account.toJSON())
        req.flash('msg', 'Confirmation message has been sent. Please check your email')
        return res.view('account/verify')
        
    } catch(e) {
        return res.serverError()
    }
}

/**
 * Verify an account
 * 
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
export async function verify(req: any, res: any) {
    try {

        let code = req.param("code")
        if (!code) {
            return res.serverError()
        }

        // get account by confirmation code
        let account = await Account.get({ confirmation_code: code })
        if (!account) {
            req.flash('error', 'Sorry, This confirmation link is invalid')
            return res.view('account/verify')
        }

        if (account.confirmed) {
            req.flash('msg', 'This account has already been confirmed')
            return res.view('account/verify')
        }

        account.set('confirmed', true)
        await account.save()
        req.flash('msg', 'Great! Your account has been confirmed')
        return res.view('account/verify')

    } catch (e) {
        console.error(e)
        return res.serverError()
    }
}