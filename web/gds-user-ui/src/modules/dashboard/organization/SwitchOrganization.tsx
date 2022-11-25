import React, { useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import TransparentLoader from 'components/Loader/TransparentLoader';
import { logUserInBff } from 'modules/auth/login/auth.service';
import { refreshNewToken } from 'utils/auth0.helper';
import { useSelector, useDispatch } from 'react-redux';
import { getAuth0User, userSelector } from 'modules/auth/login/user.slice';

import { APP_PATH } from 'utils/constants';
const SwitchOrganization: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const dispatch = useDispatch();
  const { isLoggedIn } = useSelector(userSelector);
  const navigate = useNavigate();
  const isCalled = useRef(false);

  useEffect(() => {
    const switchOrganization = async () => {
      const accessToken = await refreshNewToken();
      const { data } = await logUserInBff({
        orgid: id
      });
      console.log('data', data);

      dispatch(getAuth0User({ hasToken: accessToken }));
    };

    if (!isCalled.current) {
      switchOrganization();
      isCalled.current = true;
    }

    return () => {
      isCalled.current = false;
    };
  }, [id, dispatch]);

  useEffect(() => {
    if (isLoggedIn) {
      navigate(APP_PATH.DASHBOARD);
    }
  }, [isLoggedIn, navigate]);

  return (
    <>
      <TransparentLoader title="Switching organization ..." opacity="full" />
    </>
  );
};

export default SwitchOrganization;
