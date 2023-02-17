import React, { Suspense } from 'react';
import { Redirect, Route } from 'react-router-dom';

import PrivateRoute from './PrivateRoute';
import { lazyImport } from '@/lib/lazy-import';
import OvalLoader from '@/components/OvalLoader';

const Login = React.lazy(() => import('@/features/misc/components/account/Login'));

const { Dashboard } = lazyImport(() => import('@/features/misc'), 'Dashboard');

const { VaspsList } = lazyImport(() => import('../features/vasps'), 'VaspsList');
const { VaspDetails } = lazyImport(() => import('../features/vasps'), 'VaspDetails');

const { PageError } = lazyImport(() => import('@/features/misc'), 'PageError');
const { PageNotFound } = lazyImport(() => import('@/features/misc'), 'PageNotFound');
const { PageNotFoundAlt } = lazyImport(() => import('@/features/misc'), 'PageNotFoundAlt');
const { ServerError } = lazyImport(() => import('@/features/misc'), 'ServerError');

// root routes
const rootRoute = {
    path: '/',
    exact: true,
    children: [
        {
            path: '/',
            name: 'Project',
            component: () => <Redirect to="/dashboard" />,
            route: PrivateRoute,
            exact: true,
        },
        {
            path: '/dashboard',
            name: 'Project',
            component: Dashboard,
            route: PrivateRoute,
            exact: true,
        },
        {
            path: '/vasps/:id',
            name: 'Detail',
            component: VaspDetails,
            route: PrivateRoute,
            exact: true,
        },

        {
            path: '/vasps',
            name: 'List',
            component: VaspsList,
            route: PrivateRoute,
            exact: true,
        },
        {
            path: '/not-found',
            name: 'NotFound',
            component: PageNotFoundAlt,
            route: PrivateRoute,
        },
        {
            path: '/error',
            name: 'Error',
            component: PageError,
            route: PrivateRoute,
        },
        {
            path: '',
            name: '',
            component: () => <Redirect to="/error-404" />,
            route: Route,
        },
    ],
};

const authRoutes = [
    {
        path: '/login',
        name: 'Login',
        component: () => (
            <Suspense
                fallback={
                    <div
                        className="relative d-flex justify-content-center align-items-center"
                        style={{ height: '100vh', width: '100vw' }}>
                        <OvalLoader height={50} width={50} />
                    </div>
                }>
                <Login />
            </Suspense>
        ),
        route: Route,
    },
];

const otherPublicRoutes = [
    {
        path: '/error-404',
        name: 'Error - 404',
        component: PageNotFound,
        route: PrivateRoute,
    },
    {
        path: '/error-500',
        name: 'Error - 500',
        component: ServerError,
        route: Route,
    },
];

// flatten the list of all nested routes
const flattenRoutes = (routes) => {
    let flatRoutes = [];

    routes = routes || [];
    routes.forEach((item) => {
        flatRoutes.push(item);

        if (typeof item.children !== 'undefined') {
            flatRoutes = [...flatRoutes, ...flattenRoutes(item.children)];
        }
    });
    return flatRoutes;
};

// All routes
const authProtectedRoutes = [rootRoute];
const publicRoutes = [...authRoutes, ...otherPublicRoutes];

const authProtectedFlattenRoutes = flattenRoutes([...authProtectedRoutes]);
const publicProtectedFlattenRoutes = flattenRoutes([...publicRoutes]);

export { authProtectedFlattenRoutes, authProtectedRoutes, publicProtectedFlattenRoutes, publicRoutes };
