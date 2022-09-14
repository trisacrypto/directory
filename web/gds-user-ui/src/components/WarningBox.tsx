import { WarningIcon } from '@chakra-ui/icons';
import { StackProps, HStack, Box } from '@chakra-ui/react';

const WarningBox = ({ children, ...props }: StackProps) => (
  <HStack bg="#fff9e9" px={4} py={5} rounded="lg" fontWeight={400} {...props}>
    <WarningIcon alignSelf="start" color="#ffc12d" fontSize={{ base: 'xl', xl: '2xl' }} mt="5px" />
    <Box>{children}</Box>
  </HStack>
);

export default WarningBox;
