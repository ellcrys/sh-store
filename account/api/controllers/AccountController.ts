import { setFlagsFromString } from 'v8';
import { access } from 'fs';
import express = require('express');
import validator = require('validator');
import Account = require('../models/Account');
import ResetToken = require('../models/ResetToken');
import uuid4 = require('uuid4');
import random = require('randomstring');
import Random = require('random-js');
import bcrypt = require('bcrypt');
import sha1 = require('sha1');
import moment = require('moment');

/**
 * Sets flash error message for a field
 * 
 * @param {*} req Request object
 * @param {string} msg The error message
 * @param {string} field The field name
 * @param {*} body Response body
 */
function setFlashError(req: any, msg: string, field: string, body: any) {
    req.flash("error", msg)
    req.flash("field", field)
    req.flash("form", body)
}


/**
 * Validates res.body containing account details
 * @export
 * @param {Account} account 
 * @returns {ValidationError} 
 */
export function validateAccount(account: models.Account): ValidationError {
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

        let body = (req.body as models.Account)
        const valError = validateAccount(body)

        // send validation error
        if (valError != null) {
            setFlashError(req, valError.msg, valError.field, body)
            return res.redirect('back')
        }

        // check if an account with matching email exists
        let existingAccount = await Account.get({ email: body.email })
        if (existingAccount) {
            setFlashError(req, "Email has already been registered", "email", body)
            return res.redirect('back')
        }

        // set id, client id and secret
        body.id = uuid4()
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

        if (!body.email) {
            setFlashError(req, 'Please enter your email', 'email', body)
            return res.redirect('back')
        }

        if (!body.password) {
            setFlashError(req, 'Please enter your password', 'password', body)
            return res.redirect('back')
        }

        let account = await Account.get({ email: body.email })
        if (!account) {
            setFlashError(req, 'Your email and password combination is not valid', '', body)
            return res.redirect('back')
        }

        if (!bcrypt.compareSync(body.password, account.password)) {
            setFlashError(req, 'Your email and password combination is not valid', '', body)
            return res.redirect('back')
        }

        if (!account.confirmed) {
            req.flash('warn', 'warn')
            setFlashError(req, `Your account has not been confirmed. Please check your email for the confirmation message. <a href="/account/resend-code/${account.id}">Resend</a>`, '', body)
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

    } catch (e) {
        return res.serverError()
    }
}

/**
 * Send reset password token to email
 * 
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
export async function forgotPassword(req: any, res: any) {
    try {
        let email = req.body.email

        if (!email) {
            setFlashError(req, "Please enter your email address", "email", req.body)
            return res.view('account/forgot_password')
        }

        // get account by email
        let account = await Account.get({ email })
        if (!account) {
            req.flash('error', `No account matching the email '${req.body.email}'`)
            return res.view('account/forgot_password')
        }

        let resetToken = await ResetToken.create({ token: sha1(uuid4()), account: account.id })

        // send password reset message
        account = account.toJSON()
        account.token = resetToken.token
        let resp = await EmailService.sendPasswordResetEmail(account.email, account)
        req.flash('msg', 'A reset message has been sent to you')
        return res.view('account/forgot_password')
    } catch (e) {
        return res.serverError()
    }
}

/**
 * Reset a password page
 * 
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
export async function resetPasswordPage(req: any, res: any) {
    try {
        let resetToken = await ResetToken.get({ token: req.param("token") })
        if (!resetToken) {
            req.flash('error2', `Reset link is invalid or has expired`)
            return res.view('account/reset_password')
        }

        if (resetToken.used) {
            req.flash('error2', `Reset link is invalid or has expired`)
            return res.view('account/reset_password')
        }

        if (!moment(resetToken.created_at).add(2, "day").utc().isAfter(moment())) {
            req.flash('error2', `Reset link is invalid or has expired`)
        }

        return res.view('account/reset_password')

    } catch (e) {
        return res.serverError()
    }
}

/**
 * Reset a password
 * 
 * @export
 * @param {express.Request} req The request object
 * @param {express.Response} res The response object
 */
export async function resetPassword(req: any, res: any) {
    try {

        let resetToken = await ResetToken.get({ token: req.param("token") })
        if (!resetToken) {
            req.flash('error2', `Reset link is invalid or has expired`)
            return res.redirect(req.originalUrl)
        }

        if (resetToken.used) {
            req.flash('error2', `Reset link is invalid or has expired`)
            return res.redirect(req.originalUrl)
        }

        if (!moment(resetToken.created_at).add(2, "day").utc().isAfter(moment())) {
            req.flash('error2', `Reset link is invalid or has expired`)
            return res.redirect(req.originalUrl)
        }

        let body = req.body

        // validation
        if (validator.isEmpty(body.password)) {
            setFlashError(req, `Password is required`, 'password', body)
            return res.redirect(req.originalUrl)
        } else if (!validator.isLength(body.password, { min: 8, max: 120 })) {
            setFlashError(req, `Please enter a password between 8 and 120 characters, the trickier the better!`, 'password', body)
            return res.redirect(req.originalUrl)
        }

        let account = await Account.get({ id: resetToken.account })
        if (!account) {
            req.flash('error2', `Reset link is invalid or has expired'`)
            return res.redirect(req.originalUrl)
        }

        // ensure new password is not the same as the old one
        if (bcrypt.compareSync(body.password, account.password)) {
            setFlashError(req, `New password cannot be similar to a previously used password`, 'password', body)
            return res.redirect(req.originalUrl)
        }

        // set reset token as used
        resetToken.set('used', true)
        await resetToken.save()

        // update password
        let newPasswordHash = bcrypt.hashSync(body.password, 10)
        account.set('password', newPasswordHash)
        await account.save()

        return res.view('account/reset_password', { msg: 'Password successfully changed' })
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