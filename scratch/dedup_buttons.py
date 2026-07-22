import subprocess
import re

with open('manager/dist/assets/index-B7iM18f_.js', 'r', encoding='utf-8', errors='ignore') as f:
    text = f.read()

# Pattern to find Chatwoot button in Fz
btn_pattern = r'cw&&m\.jsxs\(m\.Fragment,\{children:\[m\.jsx\(it,\{variant:\"ghost\",className:\"rounded-none h-12 px-4 text-teal-500.*?\}\),m\.jsx\(\"div\",\{className:\"w-px bg-sidebar-border\"\}\)\]\}\),'

# Replace all occurrences of btn_pattern with empty string
text = re.sub(btn_pattern, '', text)

# Now insert it EXACTLY ONCE before target_btn
target_btn = 'm.jsx(it,{variant:"ghost",className:"rounded-none h-12 px-4 text-gray-500 hover:text-gray-300 hover:bg-gray-500/10",onClick:()=>r(t),children:m.jsx(iA,{className:"h-4 w-4"})})'
single_btn = 'cw&&m.jsxs(m.Fragment,{children:[m.jsx(it,{variant:"ghost",className:"rounded-none h-12 px-4 text-teal-500 hover:text-teal-300 hover:bg-teal-500/20 cursor-pointer transition-all duration-150",onClick:()=>cw(t),title:"Configurações do Chatwoot",children:m.jsx(CwIcon,{className:"h-4 w-4"})}),m.jsx("div",{className:"w-px bg-sidebar-border"})]}),' + target_btn

text = text.replace(target_btn, single_btn, 1)

with open('manager/dist/assets/index-B7iM18f_.js', 'w', encoding='utf-8') as f:
    f.write(text)

import subprocess
res = subprocess.run(['node', '-c', 'manager/dist/assets/index-B7iM18f_.js'], capture_output=True, text=True)
if res.returncode == 0:
    print('SUCCESS! EXACTLY 1 CHATWOOT BUTTON!')
else:
    print('FAILURE! Node stderr:\n', res.stderr)
