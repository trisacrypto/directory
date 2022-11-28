import React, { useEffect, useRef } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import TransparentLoader from 'components/Loader/TransparentLoader';
import { logUserInBff } from 'modules/auth/login/auth.service';
import { refreshNewToken } from 'utils/auth0.helper';
import { useDispatch } from 'react-redux';
// import { getAuth0User } from 'modules/auth/login/user.slice';
import { APP_PATH } from 'utils/constants';
const SwitchOrganization: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const isCalled = useRef(false);

  useEffect(() => {
    const switchOrganization = async () => {
      await refreshNewToken();
      const { data } = await logUserInBff({
        orgid: id
      });
      if (data) {
        navigate(APP_PATH.DASHBOARD);
      }
    };

    if (!isCalled.current) {
      switchOrganization();
      isCalled.current = true;
    }

    return () => {
      isCalled.current = false;
    };
  }, [id, navigate, dispatch]);

  return (
    <>
      <TransparentLoader title="Switching organization ..." opacity="full" />
    </>
  );
};

export default SwitchOrganization;
