import React from 'react'
import "/main.css"

function Login() {
    return (
        <div className="container">
            <form className="form" >
            <h1 className="form__title">Login</h1>
            <img src="./chaudrlogo.png" width="200" height="200" alt="logo"></img>
            <div className="form__message form__message--error"></div>
            <div className="form__input-group">
                <input className="form__input form__message--error" autofocus placeholder="Password"/>
                <div className="form__input-error-message"></div> 
            </div>
            <button className="form__button" type="submit">Continue</button>
            
            
            </form>            
        </div>
    )
}

export default Login
