// hook that disable button is permission is not set
import { isCurrentUser } from 'components/Collaborators/lib';
import { hasPermission } from 'utils/permission';
import { useEffect, useState } from 'react';

export const useSafeDisableIconButton = (permission: TUserPermission, condition: string) => {
  const [isDisabled, setIsDisabled] = useState(false);

  useEffect(() => {
    let once = false;
    console.log('[useSafeDisableIconButton] mount');
    if (!once) {
      once = true;
      const d = !isCurrentUser(condition) && hasPermission(permission);
      setIsDisabled(d);
    }
    return () => {
      console.log('useSafeDisableIconButton: unmount');
      once = true;
    };
  }, [permission, condition]);

  return { isDisabled };
};
