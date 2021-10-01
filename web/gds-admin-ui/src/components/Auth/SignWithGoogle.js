

import React from 'react'

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

export default SignWithGoogle
