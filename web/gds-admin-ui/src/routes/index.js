import React from 'react';
import { Redirect } from 'react-router-dom';
import { Route } from 'react-router-dom';


const Dashboard = React.lazy(() => import('../pages/app/dashboard'));
const VaspsList = React.lazy(() => import('../pages/app/lists'));


const ErrorPageNotFound = React.lazy(() => import('../pages/error/PageNotFound'));
const ServerError = React.lazy(() => import('../pages/error/ServerError'));

// root routes
const rootRoute = {
    path: '/',
    exact: true,
    component: () => <Redirect to="/dashboard" />,
    route: Route,
};

const dashboardRoutes = {
    path: '/dashboard',
    name: 'Dashboards',
    icon: 'uil-home-alt',
    header: 'Navigation',
    children: [
        {
            path: '/dashboard',
            name: 'Project',
            component: Dashboard,
            route: Route,
        },
<<<<<<< HEAD
        {
            path: '/vasps-summary/vasps',
            name: 'List',
            component: VaspsList,
            route: Route,
        }
=======
>>>>>>> feat: add dashboard page
    ],
};


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
const authProtectedRoutes = [rootRoute, dashboardRoutes];
const publicRoutes = [...otherPublicRoutes];

const authProtectedFlattenRoutes = flattenRoutes([...authProtectedRoutes]);
const publicProtectedFlattenRoutes = flattenRoutes([...publicRoutes]);

export { publicRoutes, authProtectedRoutes, authProtectedFlattenRoutes, publicProtectedFlattenRoutes };
