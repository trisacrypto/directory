import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { useSelector } from 'react-redux';

import DefaultLayout from '../layouts/Default';
import VerticalLayout from '../layouts/Vertical';

import { authProtectedFlattenRoutes, publicProtectedFlattenRoutes } from './index';
import { ModalProvider } from '@/contexts/modal';
import ResendEmail from '@/components/ResendEmail';

const Routes = (props) => {
    const { layout } = useSelector((state) => ({
        layout: state.Layout,
    }));

    const getLayout = () => {
        return VerticalLayout;
    };

    let Layout = getLayout();

    return (
        <ModalProvider>
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
            <ResendEmail />
        </ModalProvider>
    );
};

export default Routes;
