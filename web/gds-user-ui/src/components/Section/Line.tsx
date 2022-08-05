import { Box, chakra, Flex, FlexProps, StyleProps, useColorModeValue } from '@chakra-ui/react';

type Props = StyleProps &
  FlexProps & {
    children: React.ReactNode;
    title?: string;
  };

export const Line: React.FC<Props> = ({ children, title, ...rest }) => {
  return (
    <Flex>
      <Flex shrink={0}>
        <Flex rounded="md" bg={useColorModeValue('brand.500', '')} color="white"></Flex>
      </Flex>
      <Box ml={4}>
        <chakra.dt fontSize="lg" data-testid="title" fontWeight="bold" lineHeight="6" {...rest}>
          {title}
        </chakra.dt>
        <chakra.dd mt={2}>{children}</chakra.dd>
      </Box>
    </Flex>
  );
};
