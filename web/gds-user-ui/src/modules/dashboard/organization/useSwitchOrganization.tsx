import { useEffect, useState, useRef } from 'react';
import { logUserInBff, getUserCurrentOrganizationService } from 'modules/auth/login/auth.service';
import { refreshNewToken } from 'utils/auth0.helper';
import { setUserOrganization } from 'modules/auth/login/user.slice';

import { useDispatch } from 'react-redux';
// import { getAuth0User } from 'modules/auth/login/user.slice';
import { APP_STATUS_CODE } from 'utils/constants';

const useSwitchOrganization = (organizationId: string) => {
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const dispatch = useDispatch();
  const isCalled = useRef(false);
  useEffect(() => {
    const switchOrganization = async () => {
      try {
        const logged = await logUserInBff({
          orgid: organizationId
        });
        if (logged.status === APP_STATUS_CODE.OK || logged.status === APP_STATUS_CODE.NO_CONTENT) {
          const token = (await refreshNewToken()) as any;
          const user = token && (await getUserCurrentOrganizationService());
          dispatch(setUserOrganization(user?.data));
          setIsLoading(false);
        }
        // if (logged.status === APP_STATUS_CODE.NO_CONTENT) {
        //   const user = await getUserCurrentOrganizationService();
        //   dispatch(setUserOrganization(user?.data));
        //   setIsLoading(false);
        // }
      } catch (error) {
        setIsError(true);
      }
    };

    if (!isCalled.current) {
      if (organizationId) {
        switchOrganization();
        isCalled.current = true;
      } else {
        setIsLoading(false);
        setIsError(true);
      }
    }

    return () => {
      isCalled.current = false;
    };
  }, [organizationId, dispatch]);

  return { isLoading, isError, isSuccess: !isError && !isLoading };
};

export { useSwitchOrganization };
