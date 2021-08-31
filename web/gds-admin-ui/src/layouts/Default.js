// @flow
import React, { useEffect, Suspense } from 'react';

const loading = () => <div className=""></div>;

type DefaultLayoutProps = {
    layout: {
        layoutType: string,
        layoutWidth: string,
        leftSideBarTheme: string,
        leftSideBarType: string,
        showRightSidebar: boolean,
    },
    user: any,
    children?: any,
};

const DefaultLayout = (props: DefaultLayoutProps): React$Element<any> => {
    useEffect(() => {
        if (document.body) document.body.classList.add('authentication-bg');

        return () => {
            if (document.body) document.body.classList.remove('authentication-bg');
        };
    }, []);

    // get the child view which we would like to render
    const children = props.children || null;

    return (
        <>
            <Suspense fallback={loading()}>{children}</Suspense>
        </>
    );
};
export default DefaultLayout;
