// hook that disable button is permission is not set
import { isCurrentUser } from 'components/Collaborators/lib';
import { hasPermission } from 'utils/permission';
import { useEffect, useState } from 'react';

export const useSafeDisableButton = (permission: TUserPermission, condition: string) => {
  const [isDisabled, setIsDisabled] = useState(false);

  useEffect(() => {
    if (isCurrentUser(condition)) {
      setIsDisabled(true);
    }
    if (hasPermission(permission)) {
      setIsDisabled(true);
    }
  }, [permission, condition]);

  return { isDisabled };
};
