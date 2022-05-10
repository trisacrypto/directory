import { Box, Stack } from '@chakra-ui/react';

type SimpleDashboardLayout = {
  children: React.ReactNode;
};
export const SimpleDashboardLayout: React.FC<SimpleDashboardLayout> = ({ children }) => {
  return (
    <Stack px={58} py={10} bg={'#F7F8FC'} height="100%">
      <Box>{children}</Box>
    </Stack>
  );
};
