import { Box, Flex, FlexProps, StyleProps, Text, useColorModeValue } from '@chakra-ui/react';

type Props = StyleProps &
  FlexProps & {
    children: React.ReactNode;
    title?: string;
  };

export const Line: React.FC<Props> = ({ children, title }) => {
  return (
    <Flex>
      <Flex shrink={0}>
        <Flex rounded="md" bg={useColorModeValue('brand.500', '')} color="white"></Flex>
      </Flex>
      <Box ml={4}>
        <Text as="dt" fontSize="lg" data-testid="title" fontWeight="bold" lineHeight="6">
          {title}
        </Text>
        <Text as="dd" mt={2}>
          {children}
        </Text>
      </Box>
    </Flex>
  );
};
