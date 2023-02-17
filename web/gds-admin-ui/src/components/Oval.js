import React from 'react';

function Oval(props) {
    return (
        <svg width="40" height="40" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg" stroke="#fff" {...props}>
            <g fill="none" fillRule="evenodd">
                <g transform="translate(1 1)" strokeWidth="2.2">
                    <circle strokeOpacity="0.1" cx="18" cy="18" r="18" />
                    <path d="M36 18c0-9.94-8.06-18-18-18">
                        <animateTransform
                            attributeName="transform"
                            type="rotate"
                            from="0 18 18"
                            to="360 18 18"
                            dur=".7s"
                            repeatCount="indefinite"
                        />
                    </path>
                </g>
            </g>
        </svg>
    );
}

export default Oval;
