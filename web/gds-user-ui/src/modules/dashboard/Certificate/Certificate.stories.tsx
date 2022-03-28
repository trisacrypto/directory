import { Meta, Story } from '@storybook/react';
import store from 'application/store';
import { Provider } from 'react-redux';
import Certificate from './Certificate';

export default {
  title: 'modules/Certificate',
  component: Certificate,
  decorators: [
    (S) => (
      <Provider store={store}>
        <S />
      </Provider>
    )
  ]
} as Meta;

const Template: Story = (args) => <Certificate {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
