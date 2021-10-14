import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { useSelector } from 'react-redux';

import DefaultLayout from '../layouts/Default';
import VerticalLayout from '../layouts/Vertical';


import { authProtectedFlattenRoutes, publicProtectedFlattenRoutes } from './index';
import AuthProvider from '../contexts/auth/auth-provider';

const Routes = (props) => {
    const { layout } = useSelector((state) => ({
        layout: state.Layout,
        // user: state.Auth.user,
    }));

    const getLayout = () => {
        return VerticalLayout;
    };

    let Layout = getLayout();

    return (
        <AuthProvider>
            <BrowserRouter>
                <Switch>
                    <Route path={publicProtectedFlattenRoutes.map((r) => r['path'])}>
                        <DefaultLayout {...props} layout={layout}>
                            <Switch>
                                {publicProtectedFlattenRoutes.map((route, index) => {
                                    return (
                                        !route.children && (
                                            <route.route
                                                key={index}
                                                path={route.path}
                                                roles={route.roles}
                                                exact={route.exact}
                                                component={route.component}
                                            />
                                        )
                                    );
                                })}
                            </Switch>
                        </DefaultLayout>
                    </Route>

                    <Route path={authProtectedFlattenRoutes.map((r) => r['path'])}>
                        <Layout {...props} layout={layout}>
                            <Switch>
                                {authProtectedFlattenRoutes.map((route, index) => {
                                    return (
                                        !route.children && (
                                            <route.route
                                                key={index}
                                                path={route.path}
                                                roles={route.roles}
                                                exact={route.exact}
                                                component={route.component}
                                            />
                                        )
                                    );
                                })}
                            </Switch>
                        </Layout>
                    </Route>
                </Switch>
            </BrowserRouter>
        </AuthProvider>
    );
};

export default Routes;
