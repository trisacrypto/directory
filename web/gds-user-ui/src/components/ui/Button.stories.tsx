import { Button, ButtonProps } from '@chakra-ui/react';
import { Meta, Story } from '@storybook/react';

export default {
  title: 'ui/Button',
  component: Button
} as Meta;

const Template: Story<ButtonProps> = (args) => <Button {...args} />;

export const Standard = Template.bind({});
Standard.args = {
  children: 'Click me!'
};
