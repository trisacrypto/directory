import jwtDecode from 'jwt-decode';
import { useEffect, useMemo, useState } from 'react';

import useAuth from '@/contexts/auth/use-auth';

import { APICore } from '../helpers/api/apiCore';

const useUser = () => {
  const { isUserAuthenticated } = useAuth();
  const api = useMemo(() => new APICore(), []);

  const [user, setuser] = useState();

  useEffect(() => {
    if (isUserAuthenticated()) {
      const decodedUser = jwtDecode(api.getLoggedInUser()?.access_token);
      setuser(decodedUser);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [api]);

  return { user };
};

export default useUser;
