import Icon from '@chakra-ui/icon';
import { InfoIcon } from '@chakra-ui/icons';
import { Box } from '@chakra-ui/react';

const NotSavedSectionStatus: React.FC = () => {
  return (
    <Box>
      <Icon as={InfoIcon} color="#F29C36" w={7} h={7} /> (not saved)
    </Box>
  );
};

export { NotSavedSectionStatus };
