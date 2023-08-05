# -*- mode: python ; coding: utf-8 -*-


block_cipher = None

added_files = [('F:/项目/mirouter-ui/venv/static', 'static' )]

a = Analysis(['mirouter.py'],
             pathex=[],
             binaries=[],#这里可以放入你的环境所依赖的一些库，如TensorFlow或pyecharts等
             datas=added_files,#这里改成上面的数据文件列表则可引入资源文件夹和文件
             hiddenimports=[],
             hookspath=[],
             runtime_hooks=[],
             excludes=[],
             win_no_prefer_redirects=False,
             win_private_assemblies=False,
             cipher=block_cipher,
             noarchive=False)

pyz = PYZ(a.pure, a.zipped_data,
             cipher=block_cipher)
exe = EXE(pyz,
          a.scripts,
          a.binaries,
          a.zipfiles,
          a.datas,
          [],
          name='test',#         打包后生成的文件名称（可自行修改）
          debug=False,
          bootloader_ignore_signals=False,
          strip=False,
          upx=True,
          upx_exclude=[],
          runtime_tmpdir=None,
          console=True, 
          icon = 'favicon.ico')
          #上面的icon参数一般要自己加，并不会帮你生成，也可调用终端命令进行ico打包
