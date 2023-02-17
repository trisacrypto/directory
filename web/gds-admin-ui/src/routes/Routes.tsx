import React, { Suspense } from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { useSelector } from 'react-redux';

import DefaultLayout from '../layouts/Default';

import { authProtectedFlattenRoutes, publicProtectedFlattenRoutes } from './index';
import { ModalProvider } from '@/contexts/modal';
import ResendEmail from '@/components/ResendEmail';
import { lazyImport } from '@/lib/lazy-import';
import { getDirectoryLogo } from '@/utils';
import OvalLoader from '@/components/OvalLoader';

const { default: VerticalLayout } = lazyImport(
    () => import(/* webpackPreload: true */ '../layouts/Vertical'),
    'default'
);

const Routes = (props: any) => {
    const { layout } = useSelector((state: any) => ({
        layout: state.Layout,
    }));

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
                        <Suspense
                            fallback={
                                <div
                                    style={{
                                        height: '100vh',
                                        display: 'grid',
                                        placeItems: 'center',
                                        background: '#313A47',
                                    }}>
                                    <div className="d-flex flex-column gap-2">
                                        <img src={getDirectoryLogo()} alt="" height="38" />
                                        <OvalLoader title={<>&nbsp;</>} />
                                    </div>
                                </div>
                            }>
                            <VerticalLayout {...props} layout={layout}>
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
                            </VerticalLayout>
                        </Suspense>
                    </Route>
                </Switch>
            </BrowserRouter>
            <ResendEmail />
        </ModalProvider>
    );
};

export default Routes;
