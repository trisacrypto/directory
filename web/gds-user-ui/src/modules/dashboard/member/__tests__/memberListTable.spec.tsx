import { act, render, screen } from 'utils/test-utils';
import MemberTable from "../MemberTable";
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { mainnetMembersMockValue } from '../__mocks__';

function renderComponent() {
   return render(<MemberTable data={mainnetMembersMockValue} />);
}

describe('MemberTable', () => {
    beforeAll(() => {
        act(() => {
          dynamicActivate('en');
        });
      });

      it('should render', () => {
        const { container } = renderComponent();
        expect(container).toMatchSnapshot();
      })

        it('should render the correct table headers', () => {
            renderComponent();
            const header = screen.getByText(/Member Name/i);
            const joined = screen.getByText(/Joined/i);
            const lastUpdate = screen.getByText(/Last Updated/i);
            // const network = screen.getByText(/Network/i);
            const status = screen.getByText(/Status/i);
            expect(header).toBeInTheDocument();
            expect(joined).toBeInTheDocument();
            expect(lastUpdate).toBeInTheDocument();
            // expect(network).toBeInTheDocument();
            expect(status).toBeInTheDocument();
        });
});
