import React from 'react';
import { Box, Icon, Text, Stack } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

const StepLabel = (label: any, textColor: string, isActiveStep: boolean) => {
  return (
    <Stack spacing={1} width="100%">
      <Box h="1" bg={label?.color} borderRadius={'50px'} width={'100%'} />
      <Stack
        direction={{ base: 'column', md: 'row' }}
        alignItems={'center'}
        spacing={{ base: 0, md: 1 }}>
        <Box>
          <Icon
            as={label?.icon}
            sx={{
              path: {
                fill: label?.color
              },
              verticalAlign: 'middle'
            }}
          />
        </Box>
        <Text
          color={textColor}
          fontSize={'sm'}
          fontWeight={isActiveStep ? 'bold' : 'normal'}
          textAlign="center">
          2 <Trans id="Legal Person">Legal Person</Trans>
        </Text>
      </Stack>
    </Stack>
  );
};

export default StepLabel;
