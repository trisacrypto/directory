// @flow
import jwtDecode from 'jwt-decode';
import { useEffect, useState, useMemo } from 'react';

import { APICore } from '../helpers/api/apiCore';

const useUser = (): { user: any | void, ... } => {
    const api = useMemo(() => new APICore(), []);

    const [user, setuser] = useState();

    useEffect(() => {
        if (api.isUserAuthenticated()) {
            const decodedUser = jwtDecode(api.getLoggedInUser()?.access_token)
            setuser(decodedUser);
        }
    }, [api]);

    return { user };
};

export default useUser;
