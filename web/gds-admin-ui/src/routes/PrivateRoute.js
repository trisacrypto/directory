import useAuth from 'contexts/auth/use-auth';
import React from 'react';
import { Route, Redirect } from 'react-router-dom';


const PrivateRoute = ({ component: Component, ...rest }) => {
    const { isUserAuthenticated } = useAuth()

    return (
        <Route
            {...rest}
            render={(props) => {
                if (!isUserAuthenticated()) {
                    return <Redirect to={{ pathname: '/login', state: { from: props.location } }} />;
                }
                return <Component {...props} />;
            }}
        />
    );
};

export default PrivateRoute;
