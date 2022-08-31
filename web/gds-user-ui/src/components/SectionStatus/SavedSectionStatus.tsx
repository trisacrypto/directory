import Icon from '@chakra-ui/icon';
import { CheckCircleIcon } from '@chakra-ui/icons';
import { Box, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

const SavedSectionStatus: React.FC = () => {
  return (
    <Box>
      <Icon
        as={CheckCircleIcon}
        w={5}
        h={5}
        sx={{
          path: {
            fill: 'green.400'
          }
        }}
      />{' '}
      <Text as={'span'} fontSize={'sm'} pl={1}>
        <Trans id="(Saved)"> (Saved)</Trans>
      </Text>
    </Box>
  );
};

export { SavedSectionStatus };
