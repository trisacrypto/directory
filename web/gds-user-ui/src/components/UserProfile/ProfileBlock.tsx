import { ReactNode } from 'react';
import { Heading, VStack } from '@chakra-ui/react';

export const ProfileBlock = ({ title, children }: { title: ReactNode; children: ReactNode }) => {
  return (
    <VStack align="start" w="100%" spacing={5}>
      <Heading
        size="sm"
        textTransform="uppercase"
        display="flex"
        fontWeight={700}
        columnGap={4}
        alignItems="center"
        data-testid="profile_block_title">
        {title}
      </Heading>
      <VStack align="start" w="100%" spacing={4}>
        {children}
      </VStack>
    </VStack>
  );
};
