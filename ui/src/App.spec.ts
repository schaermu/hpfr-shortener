import { render, screen } from '@testing-library/svelte'
import App from './App.svelte'

test('says hpfr.ch', () => {
    render(App)
    const node = screen.queryByText('hpfr.ch');
    expect(node).not.toBeNull();
})