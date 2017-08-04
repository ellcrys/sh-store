import Bluebird = require("bluebird")
import Tempman = require("tempman")
import path = require("path")
import SparkPost = require('sparkpost')
const client = new SparkPost(sails.config.app.sparkPostKey)

/**
 * Send a message
 * 
 * @export
 * @param {{ tempFile: string, fromName: string, subject: string, to: string, data: any }} opt 
 * @returns {Promise<any>} 
 */
export function sendMsg(opt: { tempFile: string, fromName: string, subject: string, to: string, data: any }): Promise<any> {
    return new Promise((resolve, reject) => {
        try {
            Tempman.setDir(path.join(__dirname, './email_temp')).getFile(`${opt.tempFile}.njk`, opt.data, async (err, content) => {
                if (err) return reject(err);
                let data = await client.transmissions.send({
                    options: {},
                    content: {
                        from: { name: opt.fromName, email: 'no-reply@' + sails.config.app.emailDomain },
                        subject: opt.subject,
                        html: content
                    },
                    recipients: [
                        { address: opt.to }
                    ]
                })
                resolve(data)
            })
        } catch (e) {
            console.error("Error sending account confirmation message", e.message)
            reject(e)
        }
    })
}

/**
 * Send an account confirmation message
 * 
 * @export
 * @param {string} email 
 * @returns {Promise<any>} 
 */
export function sendAccountConfirmation(to: string, data: any): Promise<any> {
    return sendMsg({
        tempFile: 'account_confirmation',
        fromName: 'Ellcrys',
        subject: 'Welcome to Ellcrys, please confirm your account',
        to,
        data,
    })
}

/**
 * Send a password reset email
 * 
 * @export
 * @param {string} email 
 * @returns {Promise<any>} 
 */
export function sendPasswordResetEmail(to: string, data: any): Promise<any> {
    return sendMsg({
        tempFile: 'reset-password',
        fromName: 'Ellcrys',
        subject: 'Reset your Ellcrys password',
        to,
        data,
    })
}