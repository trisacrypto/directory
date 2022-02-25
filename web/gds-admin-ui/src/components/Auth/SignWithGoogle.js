import React from 'react'
import PropTypes from 'prop-types'

function SignWithGoogle({ loginResponse, text, clientId }) {

    React.useEffect(() => {
        if (window.google) {
            window.google.accounts.id.initialize({
                client_id: clientId,
                callback: loginResponse,
                text,
            });
            window.google.accounts.id.renderButton(
                document.getElementById("buttonDiv"),
                { theme: "filled_blue", size: "large", margin: "0 auto" }
            );
        }
    }, [clientId, loginResponse, text])

    return (
        <div style={{ display: "flex", justifyContent: "center" }}>
            <div id="buttonDiv">
            </div>
        </div>
    )
}

SignWithGoogle.propTypes = {
    text: PropTypes.oneOf(['signin_with', 'signup_with', 'continue_with']),
    clientId: PropTypes.string.isRequired,
    loginResponse: PropTypes.func.isRequired
}

export default SignWithGoogle
