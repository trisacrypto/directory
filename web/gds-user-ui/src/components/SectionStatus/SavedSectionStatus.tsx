import Icon from '@chakra-ui/icon';
import { CheckCircleIcon } from '@chakra-ui/icons';
import { Box } from '@chakra-ui/react';

const SavedSectionStatus: React.FC = () => {
  return (
    <Box>
      <Icon as={CheckCircleIcon} w={7} h={7} color="green.300" /> (saved)
    </Box>
  );
};

export { SavedSectionStatus };
