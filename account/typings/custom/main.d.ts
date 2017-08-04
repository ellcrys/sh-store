declare const models: any;
declare const sails: { 
    db: any, 
    config: any,
    getDB(): any
}

/**
 * Account represents an account
 * @interface Account
 */
interface Account {
    first_name: string,
    last_name: string,
    email: string,
    password: string,
    client_id: string,
    client_secret: string,
    confirmed: boolean,
    confirmation_code: string,
    created_at: number
}

/**
 * ValidationError represents an error about
 * an invalid res.body field.
 * @interface ValidationError
 */
interface ValidationError {
    msg: string,
    field: string
}

declare namespace EmailService {
    function sendAccountConfirmation(to: string, data: any): Promise<any>
    function sendPasswordResetEmail(to: string, data: any): Promise<any>
}