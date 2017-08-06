/**
 * Application specific configs
 */

module.exports.app = {
    
    // spark post access token
    sparkPostKey: process.env.SPARKPOST_KEY,
    
    // transaction email domain
    emailDomain: process.env.ELLDB_EMAIL_DOMAIN,
    
    // access token signing secret
    signingSecret: process.env.ELLDB_SIGNING_KEY
}