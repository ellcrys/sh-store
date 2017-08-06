/**
* Route Mappings
* (sails.config.routes)
*
* Your routes map URLs to views and controllers.
*
* If Sails receives a URL that doesn't match any of the routes below,
* it will check for matching files (images, scripts, stylesheets, etc.)
* in your assets directory.  e.g. `http://localhost:1337/images/foo.jpg`
* might match an image file: `/assets/images/foo.jpg`
*
* Finally, if those don't match either, the default 404 handler is triggered.
* See `api/responses/notFound.js` to adjust your app's 404 logic.
*
* Note: Sails doesn't ACTUALLY serve stuff from `assets`-- the default Gruntfile in Sails copies
* flat files from `assets` to `.tmp/public`.  This allows you to do things like compile LESS or
* CoffeeScript for the front-end.
*
* For more information on configuring custom routes, check out:
* http://sailsjs.org/#!/documentation/concepts/Routes/RouteTargetSyntax.html
*/
module.exports.routes = {
    '/': { view: 'index' },
    'get /signup': { view: 'account/signup' },
    'post /signup': 'AccountController.signup',
    'get /confirm': { view: 'account/confirmation' },
    'get /login': { view: 'account/login' },
    'post /login': 'AccountController.login',
    'get /account/verification/:code': 'AccountController.verify',
    'get /account/resend_code/:id': 'AccountController.resendConfirmation',
    'get /account/forgot-password': { view: 'account/forgot_password' },
    'post /account/forgot-password': 'AccountController.forgotPassword',
    'get /account/reset-password/:token': 'AccountController.resetPasswordPage',
    'post /account/reset-password/:token': 'AccountController.resetPassword',
    'post /oauth/token': 'OauthController.createToken',
    'get /oauth/authorize': 'OauthController.authorize'
};





