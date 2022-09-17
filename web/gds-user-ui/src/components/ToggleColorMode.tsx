import { IconButton, useColorMode } from '@chakra-ui/react';
import { MdModeNight, MdOutlineWbSunny } from 'react-icons/md';

function ToggleColorMode() {
  const { colorMode, toggleColorMode } = useColorMode();
  const isLight = colorMode === 'light';
  return (
    <IconButton
      aria-label={`swith to ${isLight ? 'light' : 'dark'}`}
      icon={isLight ? <MdOutlineWbSunny /> : <MdModeNight />}
      variant="outline"
      onClick={toggleColorMode}
    />
  );
}

export default ToggleColorMode;
