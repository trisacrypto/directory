import Icon from '@chakra-ui/icon';
import { CheckCircleIcon } from '@chakra-ui/icons';
import { Box } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

const SavedSectionStatus: React.FC = () => {
  return (
    <Box>
      <Icon
        as={CheckCircleIcon}
        w={7}
        h={7}
        sx={{
          path: {
            fill: 'green.400'
          }
        }}
      />{' '}
      <Trans id="(Saved)">(Saved)</Trans>
    </Box>
  );
};

export { SavedSectionStatus };
