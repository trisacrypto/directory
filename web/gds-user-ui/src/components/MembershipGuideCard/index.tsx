import React from 'react';
import { Box, Text, Stack, Button, HStack, VStack } from '@chakra-ui/react';

type MembershipGuideCardProps = {
  stepNumber: number;
  header: string;
  description: string;
  buttonText: string;
};

const MembershipGuideCard = ({
  stepNumber,
  header,
  description,
  buttonText
}: MembershipGuideCardProps) => (
  <Box textAlign="center" width="100%" maxWidth={300} minHeight="100%">
    <Stack gap={'1rem'} backgroundColor="#E5EDF1" p="1rem" height="100%">
      <VStack>
        <Text textAlign="center" fontSize="xl" data-testid="step" textTransform="capitalize">
          Step {stepNumber}
        </Text>
        <Text fontWeight="bold" data-testid="header" textTransform="capitalize">
          {header}
        </Text>
        <Text data-testid="description">{description}</Text>
      </VStack>
      <Box marginTop="auto !important">
        <Button
          variant="solid"
          size="md"
          display="inline-block"
          border="1px solid #221F1F"
          borderRadius={0}
          textTransform="capitalize">
          {buttonText}
        </Button>
      </Box>
    </Stack>
  </Box>
);

export default MembershipGuideCard;
