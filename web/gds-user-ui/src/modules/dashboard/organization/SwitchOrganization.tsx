import React, { useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import TransparentLoader from 'components/Loader/TransparentLoader';
import { logUserInBff } from 'modules/auth/login/auth.service';
import { refreshNewToken } from 'utils/auth0.helper';
import { useDispatch } from 'react-redux';
import { getAuth0User } from 'modules/auth/login/user.slice';
const SwitchOrganization: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const dispatch = useDispatch();

  const isCalled = useRef(false);

  useEffect(() => {
    const switchOrganization = async () => {
      const { data } = await logUserInBff({
        orgid: id
      });
      console.log('data', data);
      const accessToken = await refreshNewToken();
      dispatch(getAuth0User({ hasToken: accessToken }));
    };

    // const refreshNewTokenHandler = async () => {
    //   await refreshNewToken();
    // };

    if (!isCalled.current) {
      // refreshNewTokenHandler();
      switchOrganization();
      isCalled.current = true;
    }

    return () => {
      isCalled.current = false;
    };
  }, [id, dispatch]);

  return (
    <>
      <TransparentLoader title="Switching organization ..." opacity="full" />
    </>
  );
};

export default SwitchOrganization;
