import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import PrivateRoute from './PrivateRoute';

const Login = React.lazy(() => import('../pages/account/Login'));

const Dashboard = React.lazy(() => import('../pages/app/dashboard'));
const VaspsList = React.lazy(() => import('../pages/app/lists'));
const VaspsDetails = React.lazy(() => import('../pages/app/details'))


const ErrorPageNotFound = React.lazy(() => import('../pages/error/PageNotFound'));
const ErrorPageNotFoundAlt = React.lazy(() => import('../pages/error/PageNotFoundAlt'));
const ServerError = React.lazy(() => import('../pages/error/ServerError'));
const PageError = React.lazy(() => import('../pages/error/PageError'));

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
            exact: true
        },
        {
            path: '/dashboard',
            name: 'Project',
            component: Dashboard,
            route: PrivateRoute,
            exact: true
        },
        {
            path: '/vasps/:id',
            name: 'Detail',
            component: VaspsDetails,
            route: PrivateRoute,
            exact: true
        },

        {
            path: '/vasps',
            name: 'List',
            component: VaspsList,
            route: PrivateRoute,
            exact: true
        },
        {
            path: '/not-found',
            name: 'NotFound',
            component: ErrorPageNotFoundAlt,
            route: PrivateRoute
        },
        {
            path: '/error',
            name: 'Error',
            component: PageError,
            route: PrivateRoute
        },
        {
            path: '',
            name: '',
            component: () => <Redirect to="/error-404" />,
            route: Route
        },
    ],
};

const authRoutes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        route: Route,
    }
]

const otherPublicRoutes = [
    {
        path: '/error-404',
        name: 'Error - 404',
        component: ErrorPageNotFound,
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

export { publicRoutes, authProtectedRoutes, authProtectedFlattenRoutes, publicProtectedFlattenRoutes };
