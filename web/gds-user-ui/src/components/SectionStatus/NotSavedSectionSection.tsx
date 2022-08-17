import Icon from '@chakra-ui/icon';
import { InfoIcon } from '@chakra-ui/icons';
import { Box, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

const NotSavedSectionStatus: React.FC = () => {
  return (
    <Box>
      <Icon
        as={InfoIcon}
        w={5}
        h={5}
        sx={{
          path: {
            fill: '#F29C36'
          }
        }}
      />
      <Text as={'span'} fontSize={'sm'} pl={1}>
        <Trans id="(not saved)"> (not saved)</Trans>
      </Text>
    </Box>
  );
};

export { NotSavedSectionStatus };
