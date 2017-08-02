import Bluebird = require("bluebird")
import Tempman = require("tempman")
import path = require("path")
import SparkPost = require('sparkpost')
const client = new SparkPost(sails.config.app.sparkPostKey)

/**
 * Send an account confirmation message
 * 
 * @export
 * @param {string} email 
 * @returns {Promise<any>} 
 */
export function sendAccountConfirmation(to: string, data: any): Promise<any> {
    return new Promise((resolve, reject) => {
        try {
            Tempman.setDir(path.join(__dirname, './email_temp')).getFile(`account_confirmation.njk`, data, async (err, content) => {
                if (err) return reject(err);
                let data = await client.transmissions.send({
                    options: {},
                    content: {
                        from: { name: 'Ellcrys', email: 'no-reply@' + sails.config.app.emailDomain },
                        subject: 'Welcome to Ellcrys, please confirm your account',
                        html: content
                    },
                    recipients: [
                        { address: to }
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