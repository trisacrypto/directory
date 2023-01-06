import { Redirect, Route } from 'react-router-dom';

import useAuth from '@/contexts/auth/use-auth';
import { Suspense } from 'react';
import OvalLoader from '@/components/OvalLoader';

const PrivateRoute = ({ component: Component, ...rest }) => {
    const { isUserAuthenticated } = useAuth();

    return (
        <Route
            {...rest}
            render={(props) => {
                if (!isUserAuthenticated()) {
                    return <Redirect to={{ pathname: '/login', state: { from: props.location } }} />;
                }
                return (
                    <Suspense
                        fallback={
                            <div className="relative mt-3">
                                <OvalLoader />
                            </div>
                        }>
                        <Component {...props} />
                    </Suspense>
                );
            }}
        />
    );
};

export default PrivateRoute;
