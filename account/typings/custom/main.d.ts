
declare const sails: { 
    db: any, 
    config: any,
    getDB(): any
}

declare namespace models {
    /**
     * Account represents an account
     * @interface Account
     */
    interface Account {
        id: string
        first_name: string,
        last_name: string,
        email: string,
        password: string,
        confirmed: boolean,
        confirmation_code: string,
        created_at: string
    }

    /**
     * Contract represents an ellcrys contract
     */
    interface Contract {
        id: string,
        creator: string,
        name: string,
        client_id: string,
        client_secret: string,
        created_at: string
    }
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

declare namespace JSONAPIService {
    function error(status: string, detail: string, code: string, source?: string)
    function errors(errors: Array<{status: string, detail: string, code: string}>)
}