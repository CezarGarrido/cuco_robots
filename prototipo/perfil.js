

const PerfilCard = item => `
      <div class="w3-card w3-round w3-white">
        <div class="w3-container">
         <h4 class="w3-center">Meu Perfil</h4>
         <p class="w3-center"><img src="https://www.w3schools.com/w3images/avatar3.png" class="w3-circle" style="height:106px;width:106px" alt="Avatar"></p>
         <hr>
         <p><i class="fa fa-pencil fa-fw w3-margin-right w3-text-theme"></i>${item.nome}</p>
         <p><i class="fa fa-home fa-fw w3-margin-right w3-text-theme"></i>${plano.curso}</p>
         <p><i class="fa fa-birthday-cake fa-fw w3-margin-right w3-text-theme"></i>${plano.data_nascimento}</p>
        </div>
      </div>
`;
