import React from 'react';
import { Redirect } from 'react-router-dom';
import { Route } from 'react-router-dom';

const Login = React.lazy(() => import('../pages/account/Login'));

const Dashboard = React.lazy(() => import('../pages/app/dashboard'));
const VaspsList = React.lazy(() => import('../pages/app/lists'));
const VaspsDetails = React.lazy(() => import('../pages/app/details'))


const ErrorPageNotFound = React.lazy(() => import('../pages/error/PageNotFound'));
const ServerError = React.lazy(() => import('../pages/error/ServerError'));

// root routes
const rootRoute = {
    path: '/',
    exact: true,
    component: () => <Redirect to="/dashboard" />,
    route: Route,
};

const authRoutes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        route: Route,
    }
]

const dashboardRoutes = {
    path: '/dashboard',
    name: 'Dashboards',
    icon: 'uil-home-alt',
    header: 'Navigation',
    exact: true,
    children: [
        {
            path: '/dashboard',
            name: 'Project',
            component: Dashboard,
            route: Route,
        },
<<<<<<< HEAD
        {
            path: '/vasps',
            name: 'List',
            component: VaspsList,
            route: Route,
            exact: true
        }
=======
>>>>>>> feat: add dashboard page
    ],
};

const vaspsRoutes = {
    path: '/vasps',
    name: 'Vasps Summary',
    children: [
        {
            path: '/vasps/:id',
            name: 'Detail',
            component: VaspsDetails,
            route: Route,
        }
    ],
}

const appRoutes = [
    vaspsRoutes,
];

const otherPublicRoutes = [
    {
        path: '/error-404',
        name: 'Error - 404',
        component: ErrorPageNotFound,
        route: Route,
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
const authProtectedRoutes = [rootRoute, dashboardRoutes, ...appRoutes];
const publicRoutes = [...otherPublicRoutes, ...authRoutes];

const authProtectedFlattenRoutes = flattenRoutes([...authProtectedRoutes]);
const publicProtectedFlattenRoutes = flattenRoutes([...publicRoutes]);

export { publicRoutes, authProtectedRoutes, authProtectedFlattenRoutes, publicProtectedFlattenRoutes };
