import subprocess
import re

with open('manager/dist/assets/index-B7iM18f_.js', 'r', encoding='utf-8', errors='ignore') as f:
    text = f.read()

# Replace function cwL definition with the new cwL
cwL_start = text.find('function cwL({open:t,onClose:n,instance:r})')
fz_start = text.find('function Fz(')

if cwL_start != -1 and fz_start != -1:
    old_cwL_code = text[cwL_start:fz_start]
    
    # Import modal building blocks from apply_modal_design_fix
    import apply_modal_design_fix
    new_cwL_code = apply_modal_design_fix.cw_components.replace('function CwIcon(t){return m.jsx("svg",{xmlns:"http://www.w3.org/2000/svg",width:"24",height:"24",viewBox:"0 0 24 24",fill:"none",stroke:"currentColor",strokeWidth:"2",strokeLinecap:"round",strokeLinejoin:"round",...t,children:[m.jsx("path",{d:"M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"}),m.jsx("path",{d:"M8 10h.01"}),m.jsx("path",{d:"M12 10h.01"}),m.jsx("path",{d:"M16 10h.01"})].filter(Boolean)})}\n', '') + '\n'
    
    text = text[:cwL_start] + new_cwL_code + text[fz_start:]
    
    with open('manager/dist/assets/index-B7iM18f_.js', 'w', encoding='utf-8') as f:
        f.write(text)

    res = subprocess.run(['node', '-c', 'manager/dist/assets/index-B7iM18f_.js'], capture_output=True, text=True)
    if res.returncode == 0:
        print('SUCCESS! Bundle syntax verification passed!')
    else:
        print('FAILURE! Node stderr:', res.stderr)
else:
    print('ERROR: could not find cwL or Fz in bundle!')
