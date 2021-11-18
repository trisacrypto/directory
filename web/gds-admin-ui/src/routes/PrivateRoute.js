import React from 'react';
import { Route, Redirect } from 'react-router-dom';

import { APICore } from 'helpers/api/apiCore';


const PrivateRoute = ({ component: Component, ...rest }) => {
    const api = new APICore();

    return (
        <Route
            {...rest}
            render={(props) => {
                if (api.isUserAuthenticated() === false) {
                    return <Redirect to={{ pathname: '/login', state: { from: props.location } }} />;
                }
                return <Component {...props} />;
            }}
        />
    );
};

export default PrivateRoute;
